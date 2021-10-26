/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package pthi

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

type MockHECICommands struct {}

var message []byte

func (c *MockHECICommands) Init() error { return nil }
func (c *MockHECICommands) GetBufferSize() uint32 { return 4 }
func (c *MockHECICommands) SendMessage(buffer []byte, done *uint32) (bytesWritten uint32, err error) { return 4, nil }
func (c *MockHECICommands) ReceiveMessage(buffer []byte, done *uint32) (bytesRead uint32, err error) { 
	for i := 0; i < len(message) && i < len(buffer); i++ {
		buffer[i] = message[i]
	}
	return 0, nil
}
func (c *MockHECICommands) Close() {}

var pthi PTHICommand


func init() {
	pthi = PTHICommand{}
	pthi.heci = &MockHECICommands{}
}

func TestGetGUID(t *testing.T) {
	prepareMessage := GetUUIDResponse{
		Header: ResponseMessageHeader{},
		UUID: [16]uint8{0,1,2,3,4,5,6,7,8,9,11,12,13,14,15},
	} 
	message = []byte(fmt.Sprintf("%v", prepareMessage))
	result, err := pthi.GetCodeVersions()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// Hardware Tests
func TestGetGUIDH(t *testing.T) {
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
func TestGetDNSSuffix(t *testing.T) {
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
