/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package rps

import (
	"rpc/internal/amt"
	"rpc/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock the AMT Hardware
type MockCommands struct{}

var mebxDNSSuffix string

func (c MockCommands) Initialize() (bool, error) {
	return true, nil
}
func (c MockCommands) GetVersionDataFromME(key string) (string, error) { return "Version", nil }
func (c MockCommands) GetUUID() (string, error)                        { return "123-456-789", nil }
func (c MockCommands) GetControlMode() (int, error)                    { return 1, nil }
func (c MockCommands) GetOSDNSSuffix() (string, error)                 { return "osdns", nil }
func (c MockCommands) GetDNSSuffix() (string, error)                   { return mebxDNSSuffix, nil }
func (c MockCommands) GetCertificateHashes() ([]amt.CertHashEntry, error) {
	return []amt.CertHashEntry{}, nil
}
func (c MockCommands) GetRemoteAccessConnectionStatus() (amt.RemoteAccessStatus, error) {
	return amt.RemoteAccessStatus{}, nil
}
func (c MockCommands) GetLANInterfaceSettings(useWireless bool) (amt.InterfaceSettings, error) {
	return amt.InterfaceSettings{}, nil
}
func (c MockCommands) GetLocalSystemAccount() (amt.LocalSystemAccount, error) {
	return amt.LocalSystemAccount{Username: "Username", Password: "Password"}, nil
}

var p Payload

func (c MockCommands) InitiateLMS() {}

func init() {
	p = Payload{}
	p.AMT = MockCommands{}
}
func TestCreatePayload(t *testing.T) {
	mebxDNSSuffix = "mebxdns"
	result, err := p.createPayload("", "")
	assert.Equal(t, "Version", result.Version)
	assert.Equal(t, "Version", result.Build)
	assert.Equal(t, "Version", result.SKU)
	assert.Equal(t, "123-456-789", result.UUID)
	assert.Equal(t, "Username", result.Username)
	assert.Equal(t, "Password", result.Password)
	assert.Equal(t, 1, result.CurrentMode)
	assert.NotEmpty(t, result.Hostname)
	assert.Equal(t, "mebxdns", result.FQDN)
	assert.Equal(t, utils.ClientName, result.Client)
	assert.Len(t, result.CertificateHashes, 0)
	assert.NoError(t, err)
}
func TestCreatePayloadWithOSDNSSuffix(t *testing.T) {
	mebxDNSSuffix = ""
	result, err := p.createPayload("", "")
	assert.NoError(t, err)
	assert.Equal(t, "osdns", result.FQDN)
}
func TestCreatePayloadWithDNSSuffix(t *testing.T) {

	result, err := p.createPayload("vprodemo.com", "")
	assert.NoError(t, err)
	assert.Equal(t, "vprodemo.com", result.FQDN)
}
func TestCreateActivationRequestNoDNSSuffix(t *testing.T) {

	result, err := p.CreateActivationRequest("method", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "method", result.Method)
	assert.Equal(t, "key", result.APIKey)
	assert.Equal(t, "ok", result.Status)
	assert.Equal(t, "ok", result.Message)
	assert.Equal(t, utils.ProtocolVersion, result.ProtocolVersion)
	assert.Equal(t, utils.ProjectVersion, result.AppVersion)
}
func TestCreateActivationRequestWithDNSSuffix(t *testing.T) {

	result, err := p.CreateActivationRequest("method", "vprodemo.com", "")
	assert.NoError(t, err)
	assert.Equal(t, "method", result.Method)
	assert.Equal(t, "key", result.APIKey)
	assert.Equal(t, "ok", result.Status)
	assert.Equal(t, "ok", result.Message)
	assert.Equal(t, utils.ProtocolVersion, result.ProtocolVersion)
	assert.Equal(t, utils.ProjectVersion, result.AppVersion)
}

func TestCreateActivationResponse(t *testing.T) {

	result, err := p.CreateActivationResponse([]byte(""))
	assert.NoError(t, err)
	assert.Equal(t, "response", result.Method)
	assert.Equal(t, "key", result.APIKey)
	assert.Equal(t, "ok", result.Status)
	assert.Equal(t, "ok", result.Message)
	assert.Equal(t, utils.ProtocolVersion, result.ProtocolVersion)
	assert.Equal(t, utils.ProjectVersion, result.AppVersion)

}
