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
func (c MockPTHICommands) GetCodeVersions() (pthi.GetCodeVersionsResponse, error) {
	return pthi.GetCodeVersionsResponse{
		CodeVersion: pthi.CodeVersions{
			BiosVersion: [65]uint8{84, 101, 115, 116},
			VersionsCount: 0,
			Versions: [50]pthi.AMTVersionType{},
		},
	}, nil
}
func (c MockPTHICommands) GetUUID() (uuid string, err error) { return "\xd2?\x11\x1c%3\x94E\xa2rT\xb2\x03\x8b\xeb\a" , nil }
func (c MockPTHICommands) GetControlMode() (state int, err error) {return 0, nil}
func (c MockPTHICommands) GetDNSSuffix() (suffix string, err error) {return "Test", nil}
func (c MockPTHICommands) GetCertificateHashes() (hashEntryList []pthi.CertHashEntry, err error) {
	return []pthi.CertHashEntry{{
		CertificateHash: [64]uint8{84, 101, 115, 116},
		Name:            pthi.AMTANSIString{Length: 4, Buffer: [1000]uint8{84, 101, 115, 116}},
		HashAlgorithm:   2,
		IsActive:        1,
		IsDefault:       1,
	}}, nil
}
func (c MockPTHICommands) GetRemoteAccessConnectionStatus() (RAStatus pthi.GetRemoteAccessConnectionStatusResponse, err error) {
	return pthi.GetRemoteAccessConnectionStatusResponse{
		NetworkStatus: 2,
		RemoteStatus: 0,
		RemoteTrigger: 0,
		MPSHostname: pthi.AMTANSIString{Length: 4, Buffer: [1000]uint8{84, 101, 115, 116}},
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
	return pthi.GetLocalSystemAccountResponse{
		Status: 0,
		Account: pthi.LocalSystemAccount{
			Username: [33]uint8{84, 101, 115, 116},
			Password: [33]uint8{84, 101, 115, 116},
		},
	}, nil
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
