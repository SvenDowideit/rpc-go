/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package amt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockPTHICommands struct {}

func (c MockPTHICommands) GetUUID() (string, error) { return "\xd2?\x11\x1c%3\x94E\xa2rT\xb2\x03\x8b\xeb\a" , nil }

func init() {
	//amt := Command{}
	//amt.PTHI = MockPTHICommands{}
}

// Mocked Tests
func MockedTestGetGUID(t *testing.T) {
	amt := Command{}
	result, err := amt.GetUUID()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// Hardware Tests
func TestGetGUID(t *testing.T) {
	amt := Command{}
	result, err := amt.GetUUID()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetControlmode(t *testing.T) {
	amt := Command{}
	result, err := amt.GetControlMode()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetDNSSuffix(t *testing.T) {
	amt := Command{}
	result, err := amt.GetDNSSuffix()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetCertificateHashes(t *testing.T) {
	amt := Command{}
	result, err := amt.GetCertificateHashes()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetRemoteAccessConnectionStatus(t *testing.T) {
	amt := Command{}
	result, err := amt.GetRemoteAccessConnectionStatus()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetLANInterfaceSettingsTrue(t *testing.T) {
	amt := Command{}
	result, err := amt.GetLANInterfaceSettings(true)
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetLANInterfaceSettingsFalse(t *testing.T) {
	amt := Command{}
	result, err := amt.GetLANInterfaceSettings(false)
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetLocalSystemAccount(t *testing.T) {
	amt := Command{}
	result, err := amt.GetLocalSystemAccount()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}
