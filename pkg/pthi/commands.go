/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package pthi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"rpc/pkg/heci"
)

type PTHICommand struct {
	heci heci.Heci
}

type PTHIInfoCommands interface {
	NewPTHICommand() (PTHIInfoCommands, error)
	Close()
	Call(command []byte, commandSize uint32) (result []byte, err error)
	GetCodeVersions() (GetCodeVersionsResponse, error)
	GetUUID() (uuid string, err error)
	GetControlMode() (state int, err error)
	GetDNSSuffix() (suffix string, err error)
	GetCertificateHashes() (hashEntryList []CertHashEntry, err error)
	GetRemoteAccessConnectionStatus() (RAStatus GetRemoteAccessConnectionStatusResponse, err error)
	GetLANInterfaceSettings(useWireless bool) (LANInterface GetLANInterfaceSettingsResponse, err error)
	GetLocalSystemAccount() (localAccount GetLocalSystemAccountResponse, err error)
}

func (pthi PTHICommand) NewPTHICommand() (PTHIInfoCommands, error) {
	heci := heci.Heci{}

	err := heci.Init()
	if err != nil {
		emptyCommand := PTHICommand{}
		return emptyCommand, err
	}
	return PTHICommand{
		heci: heci,
	}, nil
}

func (pthi PTHICommand) Close() {
	pthi.heci.Close()
}

func (pthi PTHICommand) Call(command []byte, commandSize uint32) (result []byte, err error) {
	size := pthi.heci.GetBufferSize()

	bytesWritten, err := pthi.heci.SendMessage(command, &commandSize)
	if err != nil {
		return nil, err
	}
	if bytesWritten != uint32(len(command)) {
		return nil, errors.New("amt internal error")
	}
	readBuffer := make([]byte, size)
	bytesRead, err := pthi.heci.ReceiveMessage(readBuffer, &size)
	if err != nil {
		return nil, err
	}

	if bytesRead == 0 {
		return nil, errors.New("empty response from AMT")
	}
	return readBuffer, nil
}

func CreateRequestHeader(command uint32, length uint32) MessageHeader {
	return MessageHeader{
		Version: Version{
			MajorNumber: 1,
			MinorNumber: 1,
		},
		Reserved: 0,
		Command: CommandFormat{
			val: command,
		},
		Length: length,
	}
}

func (pthi PTHICommand) GetCodeVersions() (GetCodeVersionsResponse, error) {
	commandSize := (uint32)(12)
	command := GetUUIDRequest{
		Header: CreateRequestHeader(CODE_VERSIONS_REQUEST, 0),
	}
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.LittleEndian, command)
	result, err := pthi.Call(bin_buf.Bytes(), commandSize)
	if err != nil {
		return GetCodeVersionsResponse{}, err
	}
	buf2 := bytes.NewBuffer(result)
	response := GetCodeVersionsResponse{
		Header: readHeaderResponse(buf2),
	}
	binary.Read(buf2, binary.LittleEndian, &response.CodeVersion)
	// binary.Read(buf2, binary.LittleEndian, &response.CodeVersion.BiosVersion)
	// binary.Read(buf2, binary.LittleEndian, &response.CodeVersion.VersionsCount)
	// binary.Read(buf2, binary.LittleEndian, &response.CodeVersion.Versions)

	return response, nil
}

func (pthi PTHICommand) GetUUID() (uuid string, err error) {
	commandSize := (uint32)(12) //(uint32)(unsafe.Sizeof(GetUUIDRequest{}))
	command := GetUUIDRequest{
		Header: CreateRequestHeader(GET_UUID_REQUEST, 0),
	}
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.LittleEndian, command)
	result, err := pthi.Call(bin_buf.Bytes(), commandSize)
	if err != nil {
		return "", err
	}
	buf2 := bytes.NewBuffer(result)
	response := GetUUIDResponse{
		Header: readHeaderResponse(buf2),
	}

	binary.Read(buf2, binary.LittleEndian, &response.UUID)

	return string(([]byte)(response.UUID[:])), nil
}

func (pthi PTHICommand) GetControlMode() (state int, err error) {
	commandSize := (uint32)(12)
	command := GetControlModeRequest{
		Header: CreateRequestHeader(GET_CONTROL_MODE_REQUEST, 0),
	}
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.LittleEndian, command)
	result, err := pthi.Call(bin_buf.Bytes(), commandSize)
	if err != nil {
		return -1, err
	}
	buf2 := bytes.NewBuffer(result)
	response := GetControlModeResponse{
		Header: readHeaderResponse(buf2),
	}

	binary.Read(buf2, binary.LittleEndian, &response.State)

	return response.State, nil
}

func readHeaderResponse(header *bytes.Buffer) ResponseMessageHeader {
	response := ResponseMessageHeader{}

	binary.Read(header, binary.LittleEndian, &response.Header.Version.MajorNumber)
	binary.Read(header, binary.LittleEndian, &response.Header.Version.MinorNumber)
	binary.Read(header, binary.LittleEndian, &response.Header.Reserved)
	binary.Read(header, binary.LittleEndian, &response.Header.Command.val)
	// binary.Read(header, binary.LittleEndian, &response.Header.Header.Command.fields)
	binary.Read(header, binary.LittleEndian, &response.Header.Length)
	binary.Read(header, binary.LittleEndian, &response.Status)

	return response
}

func (pthi PTHICommand) GetDNSSuffix() (suffix string, err error) {
	commandSize := (uint32)(12)
	command := GetPKIFQDNSuffixRequest{
		Header: CreateRequestHeader(GET_PKI_FQDN_SUFFIX_REQUEST, 0),
	}
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.LittleEndian, command)
	result, err := pthi.Call(bin_buf.Bytes(), commandSize)
	if err != nil {
		return "", err
	}
	buf2 := bytes.NewBuffer(result)
	response := GetPKIFQDNSuffixResponse{
		Header: readHeaderResponse(buf2),
	}

	binary.Read(buf2, binary.LittleEndian, &response.Suffix.Length)
	binary.Read(buf2, binary.LittleEndian, &response.Suffix.Buffer)

	if int(response.Suffix.Length) > 0 {
		return string(response.Suffix.Buffer[:response.Suffix.Length]), nil
	}

	return "", nil
}

func (pthi PTHICommand) GetCertificateHashes() (hashEntryList []CertHashEntry, err error) {
	// Enumerate a list of hash handles to request from
	enumerateCommandSize := (uint32)(12)
	enumerateCommand := GetHashHandlesRequest{
		Header: CreateRequestHeader(ENUMERATE_HASH_HANDLES_REQUEST, 0),
	}
	var EnumerateBin_buf bytes.Buffer
	binary.Write(&EnumerateBin_buf, binary.LittleEndian, enumerateCommand)
	enumerateResult, err := pthi.Call(EnumerateBin_buf.Bytes(), enumerateCommandSize)
	if err != nil {
		emptyHashList := []CertHashEntry{}
		return emptyHashList, err
	}
	enumerateBuf2 := bytes.NewBuffer(enumerateResult)
	enumerateResponse := GetHashHandlesResponse{
		Header: readHeaderResponse(enumerateBuf2),
	}

	binary.Read(enumerateBuf2, binary.LittleEndian, &enumerateResponse.HashHandles.Length)
	binary.Read(enumerateBuf2, binary.LittleEndian, &enumerateResponse.HashHandles.Handles)

	// Request from the enumerated list and return cert hashes
	for i := 0; i < int(enumerateResponse.HashHandles.Length); i++ {
		commandSize := (uint32)(16)
		command := GetCertHashEntryRequest{
			Header:     CreateRequestHeader(GET_CERTHASH_ENTRY_REQUEST, 4),
			HashHandle: enumerateResponse.HashHandles.Handles[i],
		}
		var bin_buf bytes.Buffer
		binary.Write(&bin_buf, binary.LittleEndian, command)
		result, err := pthi.Call(bin_buf.Bytes(), commandSize)
		if err != nil {
			emptyHashList := []CertHashEntry{}
			return emptyHashList, err
		}
		buf2 := bytes.NewBuffer(result)
		response := GetCertHashEntryResponse{
			Header: readHeaderResponse(buf2),
		}

		binary.Read(buf2, binary.LittleEndian, &response.Hash.IsDefault)
		binary.Read(buf2, binary.LittleEndian, &response.Hash.IsActive)
		binary.Read(buf2, binary.LittleEndian, &response.Hash.CertificateHash)
		binary.Read(buf2, binary.LittleEndian, &response.Hash.HashAlgorithm)
		binary.Read(buf2, binary.LittleEndian, &response.Hash.Name.Length)
		binary.Read(buf2, binary.LittleEndian, &response.Hash.Name.Buffer)

		hashEntryList = append(hashEntryList, response.Hash)
	}

	return hashEntryList, nil
}

func (pthi PTHICommand) GetRemoteAccessConnectionStatus() (RAStatus GetRemoteAccessConnectionStatusResponse, err error) {
	commandSize := (uint32)(12)
	command := GetRemoteAccessConnectionStatusRequest{
		Header: CreateRequestHeader(GET_REMOTE_ACCESS_CONNECTION_STATUS_REQUEST, 0),
	}
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.LittleEndian, command)
	result, err := pthi.Call(bin_buf.Bytes(), commandSize)
	if err != nil {
		emptyResponse := GetRemoteAccessConnectionStatusResponse{}
		return emptyResponse, err
	}
	buf2 := bytes.NewBuffer(result)
	response := GetRemoteAccessConnectionStatusResponse{
		Header: readHeaderResponse(buf2),
	}

	binary.Read(buf2, binary.LittleEndian, &response.NetworkStatus)
	binary.Read(buf2, binary.LittleEndian, &response.RemoteStatus)
	binary.Read(buf2, binary.LittleEndian, &response.RemoteTrigger)
	binary.Read(buf2, binary.LittleEndian, &response.MPSHostname.Length)
	binary.Read(buf2, binary.LittleEndian, &response.MPSHostname.Buffer)

	return response, nil
}

func (pthi PTHICommand) GetLANInterfaceSettings(useWireless bool) (LANInterface GetLANInterfaceSettingsResponse, err error) {
	commandSize := (uint32)(16)
	command := GetLANInterfaceSettingsRequest{
		Header:         CreateRequestHeader(GET_LAN_INTERFACE_SETTINGS_REQUEST, 4),
		InterfaceIndex: 0,
	}
	if useWireless {
		command.InterfaceIndex = 1
	}
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.LittleEndian, command)
	result, err := pthi.Call(bin_buf.Bytes(), commandSize)
	if err != nil {
		emptySettings := GetLANInterfaceSettingsResponse{}
		return emptySettings, err
	}
	buf2 := bytes.NewBuffer(result)
	response := GetLANInterfaceSettingsResponse{
		Header: readHeaderResponse(buf2),
	}

	binary.Read(buf2, binary.LittleEndian, &response.Enabled)
	binary.Read(buf2, binary.LittleEndian, &response.Ipv4Address)
	binary.Read(buf2, binary.LittleEndian, &response.DhcpEnabled)
	binary.Read(buf2, binary.LittleEndian, &response.DhcpIpMode)
	binary.Read(buf2, binary.LittleEndian, &response.LinkStatus)
	binary.Read(buf2, binary.LittleEndian, &response.MacAddress)

	return response, nil
}

func (pthi PTHICommand) GetLocalSystemAccount() (localAccount GetLocalSystemAccountResponse, err error) {
	commandSize := (uint32)(52)
	command := GetLocalSystemAccountRequest{
		Header: CreateRequestHeader(GET_LOCAL_SYSTEM_ACCOUNT_REQUEST, 40),
	}
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.LittleEndian, command)
	result, err := pthi.Call(bin_buf.Bytes(), commandSize)
	if err != nil {
		emptyAccount := GetLocalSystemAccountResponse{}
		return emptyAccount, err
	}
	buf2 := bytes.NewBuffer(result)
	response := GetLocalSystemAccountResponse{
		Header: readHeaderResponse(buf2),
	}

	binary.Read(buf2, binary.LittleEndian, &response.Account.Username)
	binary.Read(buf2, binary.LittleEndian, &response.Account.Password)

	return response, nil
}
