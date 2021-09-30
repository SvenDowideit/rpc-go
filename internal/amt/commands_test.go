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
	result, err := amt.GetUUIDV2()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetControlmode(t *testing.T) {
	amt := Command{}
	result, err := amt.GetControlModeV2()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetDNSSuffix(t *testing.T) {
	amt := Command{}
	result, err := amt.GetDNSSuffixV2()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetCertificateHashesV2(t *testing.T) {
	amt := Command{}
	result, err := amt.GetCertificateHashesV2()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetRemoteAccessConnectionStatusV2(t *testing.T) {
	amt := Command{}
	result, err := amt.GetRemoteAccessConnectionStatusV2()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetLANInterfaceSettingsV2True(t *testing.T) {
	amt := Command{}
	result, err := amt.GetLANInterfaceSettingsV2(true)
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetLANInterfaceSettingsV2False(t *testing.T) {
	amt := Command{}
	result, err := amt.GetLANInterfaceSettingsV2(false)
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetLocalSystemAccountV2(t *testing.T) {
	amt := Command{}
	result, err := amt.GetLocalSystemAccountV2()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

// sudo /usr/local/go/bin/go test -timeout 30s -run ^TestGetCertificateHashesV2$ rpc/internal/amt