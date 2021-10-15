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
	pthi := NewPTHICommand()
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetUUID()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetCodeVersions(t *testing.T) {
	pthi := NewPTHICommand()
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetCodeVersions()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

}
func TestGetDNSSuffixV2(t *testing.T) {
	pthi := NewPTHICommand()
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetDNSSuffix()

	assert.NoError(t, err)
	assert.Empty(t, result)

}

func TestGetCertificateHashes(t *testing.T) {
	pthi := NewPTHICommand()
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetCertificateHashes()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetRemoteAccessConnectionStatus(t *testing.T) {
	pthi := NewPTHICommand()
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetRemoteAccessConnectionStatus()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetLANInterfaceSettingsTrue(t *testing.T) {
	pthi := NewPTHICommand()
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetLANInterfaceSettings(true)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetLANInterfaceSettingsFalse(t *testing.T) {
	pthi := NewPTHICommand()
	err := pthi.Open()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetLANInterfaceSettings(false)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetLocalSystemAccount(t *testing.T) {
	pthi := NewPTHICommand()
	err := pthi.heci.Init()
	defer pthi.Close()
	assert.NoError(t, err)
	result, err := pthi.GetLocalSystemAccount()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
