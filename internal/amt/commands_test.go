/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package amt

import (
	"rpc/pkg/pthi"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockPTHICommands struct {}

func (c MockPTHICommands) Open() error { return nil }
func (c MockPTHICommands) Close() {}
func (c MockPTHICommands) Call(command []byte, commandSize uint32) (result []byte, err error) {return nil, nil}
func (c MockPTHICommands) GetCodeVersions() (pthi.GetCodeVersionsResponse, error) {return pthi.GetCodeVersionsResponse{}, nil}
func (c MockPTHICommands) GetUUID() (uuid string, err error) { return "\xd2?\x11\x1c%3\x94E\xa2rT\xb2\x03\x8b\xeb\a" , nil }
func (c MockPTHICommands) GetControlMode() (state int, err error) {return 0, nil}
func (c MockPTHICommands) GetDNSSuffix() (suffix string, err error) {return "test", nil}
func (c MockPTHICommands) GetCertificateHashes() (hashEntryList []pthi.CertHashEntry, err error) {
	var buf [1000]uint8
	buf[0] = 'A'
	hash := []pthi.CertHashEntry{}
	hash = append(hash, pthi.CertHashEntry{
		CertificateHash: [64]uint8{231,104,86,52,239,172,246,154,206,147,154,107,37,91,123,79,171,239,66,147,91,80,162,101,172,181,203,96,39,228,78,112,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
		Name: pthi.AMTANSIString{Length: 1, Buffer: buf},
		HashAlgorithm: 2,
		IsActive: 1,
		IsDefault: 1,
	})
	return hash, nil
}
func (c MockPTHICommands) GetRemoteAccessConnectionStatus() (RAStatus pthi.GetRemoteAccessConnectionStatusResponse, err error) {
	var buf [1000]uint8
	buf[0] = 'A'
	return pthi.GetRemoteAccessConnectionStatusResponse{
		NetworkStatus: 2,
		RemoteStatus: 0,
		RemoteTrigger: 0,
		MPSHostname: pthi.AMTANSIString{Length: 1, Buffer: buf},

	}, nil
}
func (c MockPTHICommands) GetLANInterfaceSettings(useWireless bool) (LANInterface pthi.GetLANInterfaceSettingsResponse, err error) {
	if useWireless {
		return pthi.GetLANInterfaceSettingsResponse{
			Enabled: 0,
			Ipv4Address: 0,
			DhcpEnabled: 1,
			DhcpIpMode: 0,
			LinkStatus: 0,
			MacAddress: [6]uint8{0,0,0,0,0,0},
		}, nil

	} else {
		return pthi.GetLANInterfaceSettingsResponse{
			Enabled: 1,
			Ipv4Address: 0,
			DhcpEnabled: 1,
			DhcpIpMode: 2,
			LinkStatus: 1,
			MacAddress: [6]uint8{7,7,7,7,7,7},
		}, nil
	}
}
func (c MockPTHICommands) GetLocalSystemAccount() (localAccount pthi.GetLocalSystemAccountResponse, err error) {
	return pthi.GetLocalSystemAccountResponse{}, nil
}

var amt Command

func init() {
	amt = Command{}
	amt.PTHI = MockPTHICommands{}
}

// Tests
func TestGetGUID(t *testing.T) {
	result, err := amt.GetUUID()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetControlmode(t *testing.T) {
	result, err := amt.GetControlMode()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetDNSSuffix(t *testing.T) {
	result, err := amt.GetDNSSuffix()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetCertificateHashes(t *testing.T) {
	result, err := amt.GetCertificateHashes()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetRemoteAccessConnectionStatus(t *testing.T) {
	result, err := amt.GetRemoteAccessConnectionStatus()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetLANInterfaceSettingsTrue(t *testing.T) {
	result, err := amt.GetLANInterfaceSettings(true)
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetLANInterfaceSettingsFalse(t *testing.T) {
	result, err := amt.GetLANInterfaceSettings(false)
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetLocalSystemAccount(t *testing.T) {
	result, err := amt.GetLocalSystemAccount()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}
