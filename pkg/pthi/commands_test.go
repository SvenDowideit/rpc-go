/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package pthi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGUID(t *testing.T) {
	pthi := PTHICommand{}
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetUUID()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetDNSSuffixV2(t *testing.T) {
	pthi := PTHICommand{}
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetDNSSuffix()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

}

func TestGetCertificateHashes(t *testing.T) {
	pthi := PTHICommand{}
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetCertificateHashes()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetRemoteAccessConnectionStatus(t *testing.T) {
	pthi := PTHICommand{}
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetRemoteAccessConnectionStatus()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetLANInterfaceSettingsTrue(t *testing.T) {
	pthi := PTHICommand{}
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetLANInterfaceSettings(true)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetLANInterfaceSettingsFalse(t *testing.T) {
	pthi := PTHICommand{}
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetLANInterfaceSettings(false)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// sudo /usr/local/go/bin/go test -timeout 30s -run ^TestGetCertificateHashes$ rpc/pkg/pthi
