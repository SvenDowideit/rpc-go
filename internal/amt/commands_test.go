/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package amt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// sudo /usr/local/go/bin/go test -timeout 30s -run ^TestGetCertificateHashes$ rpc/internal/amt
