//go:build amt
// +build amt

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
func (c MockPTHICommands) GetCertificateHashes() (hashEntryList []pthi.CertHashEntry, err error) {return []pthi.CertHashEntry{}, nil}
func (c MockPTHICommands) GetRemoteAccessConnectionStatus() (RAStatus pthi.GetRemoteAccessConnectionStatusResponse, err error) {
	return pthi.GetRemoteAccessConnectionStatusResponse{}, nil
}
func (c MockPTHICommands) GetLANInterfaceSettings(useWireless bool) (LANInterface pthi.GetLANInterfaceSettingsResponse, err error) {
	return pthi.GetLANInterfaceSettingsResponse{}, nil
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
