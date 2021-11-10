/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package rpc

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"rpc/internal/amt"
	"rpc/pkg/utils"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Flags holds data received from the command line
type Flags struct {
	commandLineArgs       []string
	URL                   string
	DNS                   string
	Hostname              string
	Proxy                 string
	Command               string
	Profile               string
	SkipCertCheck         bool
	Verbose               bool
	SyncClock             bool
	amtInfoCommand        *flag.FlagSet
	amtActivateCommand    *flag.FlagSet
	amtDeactivateCommand  *flag.FlagSet
	amtMaintenanceCommand *flag.FlagSet
}

// Stores AMTMessage for amtinfo
type RASMessage struct {
	Network string `json:"network"`
	RemoteStatus string `json:"remoteStatus"`
	Trigger string `json:"trigger"`
	MPSHostname string `json:"mpsHostname"`
}
type LANInterfaceMessage struct {
	DHCPEnabled string `json:"dhcpEnabled"`
	DHCPMode string `json:"dhcpMode"`
	LinkStatus string `json:"linkStatus"`
	IPAddress string `json:"ipAdress"`
	MACAddress string `json:"macAddress"`
}
type CertificateHashMessage struct {
	IsDefault bool `json:"isDefault"`
	IsActive bool `json:"isActive"`
	Algorithm string `json:"algorithm"`
	Hash string `json:"hash"`
}
// Pointers are used so that omitempty will ignore 
// fields the user does not want while still printing 
// empty values if that is what AMT returns.
type AMTInfoMessage struct {
	Version *string `json:"version,omitempty"`
	BuildNumber *string `json:"buildNumber,omitempty"`
	SKU *string `json:"sku,omitempty"`
	UUID *string `json:"uuid,omitempty"`
	ControlMode *string `json:"controlMode,omitempty"`
	DNSSuffix *string `json:"dnsSuffix,omitempty"`
	DNSSuffixOS *string `json:"dnsSuffixOs,omitempty"`
	Hostname *string `json:"hostname,omitempty"`
	RAS *RASMessage `json:"ras,omitempty"`
	Wired *LANInterfaceMessage `json:"wired,omitempty"`
	Wireless *LANInterfaceMessage `json:"wireless,omitempty"`
	CertificateHashes *[]CertificateHashMessage `json:"certHashes,omitempty"`
}

func NewFlags(args []string) *Flags {
	flags := &Flags{}
	flags.commandLineArgs = args
	flags.amtInfoCommand = flag.NewFlagSet("amtinfo", flag.ExitOnError)
	flags.amtActivateCommand = flag.NewFlagSet("activate", flag.ExitOnError)
	flags.amtDeactivateCommand = flag.NewFlagSet("deactivate", flag.ExitOnError)
	flags.amtMaintenanceCommand = flag.NewFlagSet("maintenance", flag.ExitOnError)
	flags.setupCommonFlags()
	return flags
}

// ParseFlags is used for understanding the command line flags
func (f *Flags) ParseFlags() (string, bool) {

	if len(f.commandLineArgs) > 1 {
		switch f.commandLineArgs[1] {
		case "amtinfo":
			f.handleAMTInfo(f.amtInfoCommand)
			return "amtinfo", false //we want to exit the program
		case "activate":
			success := f.handleActivateCommand()
			return "activate", success
		case "maintenance":
			success := f.handleMaintenanceCommand()
			return "maintenance", success
		case "deactivate":
			success := f.handleDeactivateCommand()
			return "deactivate", success
		case "version":
			println(strings.ToUpper(utils.ProjectName))
			println("Protocol " + utils.ProtocolVersion)
			return "version", false
		default:
			f.printUsage()
			return "", false
		}
	}
	f.printUsage()
	return "", false

}
func (f *Flags) printUsage() string {
	usage := "\nRemote Provisioning Client (RPC) - used for activation, deactivation, and status of AMT\n\n"
	usage = usage + "Usage: rpc COMMAND [OPTIONS]\n\n"
	usage = usage + "Supported Commands:\n"
	usage = usage + "  activate    Activate this device with a specified profile\n"
	usage = usage + "              Example: ./rpc activate -u wss://server/activate --profile acmprofile\n"
	usage = usage + "  deactivate  Deactivates this device. AMT password is required\n"
	usage = usage + "              Example: ./rpc deactivate -u wss://server/activate\n"
	usage = usage + "  maintenance Maintain this device.\n"
	usage = usage + "              Example: ./rpc maintenance -u wss://server/activate\n"
	usage = usage + "  amtinfo     Displays information about AMT status and configuration\n"
	usage = usage + "              Example: ./rpc amtinfo\n"
	usage = usage + "  version     Displays the current version of RPC and the RPC Protocol version\n"
	usage = usage + "              Example: ./rpc version\n"
	usage = usage + "\nRun 'rpc COMMAND' for more information on a command.\n"
	fmt.Println(usage)
	return usage
}

func (f *Flags) setupCommonFlags() {
	for _, fs := range []*flag.FlagSet{f.amtActivateCommand, f.amtDeactivateCommand, f.amtMaintenanceCommand} {
		fs.StringVar(&f.URL, "u", "", "websocket address of server to activate against") //required
		fs.BoolVar(&f.SkipCertCheck, "n", false, "skip websocket server certificate verification")
		fs.StringVar(&f.Proxy, "p", "", "proxy address and port")
		fs.BoolVar(&f.Verbose, "v", false, "verbose output")
	}
}
func (f *Flags) handleMaintenanceCommand() bool {
	passwordPtr := f.amtMaintenanceCommand.String("password", "", "AMT password")
	f.amtMaintenanceCommand.BoolVar(&f.SyncClock, "c", false, "sync AMT clock")
	if len(f.commandLineArgs) == 2 {
		f.amtMaintenanceCommand.PrintDefaults()
		return false
	}
	f.amtMaintenanceCommand.Parse(f.commandLineArgs[2:])
	if f.amtMaintenanceCommand.Parsed() {
		if f.URL == "" {
			fmt.Println("-u flag is required and cannot be empty")
			f.amtActivateCommand.Usage()
			return false
		}
		if *passwordPtr == "" {
			fmt.Println("Please enter AMT Password: ")
			var password string
			// Taking input from user
			_, err := fmt.Scanln(&password)
			if password == "" || err != nil {
				return false
			}
			*passwordPtr = password
		}
	}
	f.Command = "maintenance --synctime --password " + *passwordPtr
	return true
}

func (f *Flags) lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
func (f *Flags) lookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		parsedVal, err := strconv.ParseBool(val)
		if err != nil {
			log.Error(err)
			return false
		}
		return parsedVal
	}
	return defaultVal
}

func (f *Flags) handleActivateCommand() bool {
	f.amtActivateCommand.StringVar(&f.DNS, "d", f.lookupEnvOrString("DNS_SUFFIX", ""), "dns suffix override")
	f.amtActivateCommand.StringVar(&f.Hostname, "h", f.lookupEnvOrString("HOSTNAME", ""), "hostname override")
	f.amtActivateCommand.StringVar(&f.Profile, "profile", f.lookupEnvOrString("PROFILE", ""), "name of the profile to use")
	if len(f.commandLineArgs) == 2 {
		f.amtActivateCommand.PrintDefaults()
		return false
	}
	f.amtActivateCommand.Parse(f.commandLineArgs[2:])

	if f.amtActivateCommand.Parsed() {
		if f.URL == "" {
			fmt.Println("-u flag is required and cannot be empty")
			f.amtActivateCommand.Usage()
			return false
		}
		if f.Profile == "" {
			fmt.Println("-profile flag is required and cannot be empty")
			f.amtActivateCommand.Usage()
			return false
		}
	}
	f.Command = "activate --profile " + f.Profile
	return true
}
func (f *Flags) handleDeactivateCommand() bool {
	passwordPtr := f.amtDeactivateCommand.String("password", f.lookupEnvOrString("AMT_PASSWORD", ""), "AMT password")
	forcePtr := f.amtDeactivateCommand.Bool("f", false, "force deactivate even if device is not registered with a server")

	if len(f.commandLineArgs) == 2 {
		f.amtDeactivateCommand.PrintDefaults()
		return false
	}
	f.amtDeactivateCommand.Parse(f.commandLineArgs[2:])

	if f.amtDeactivateCommand.Parsed() {
		if f.URL == "" {
			fmt.Println("-u flag is required and cannot be empty")
			f.amtDeactivateCommand.Usage()
			return false
		}
		if *passwordPtr == "" {
			fmt.Println("Please enter AMT Password: ")
			var password string
			// Taking input from user
			_, err := fmt.Scanln(&password)
			if password == "" || err != nil {
				return false
			}
			*passwordPtr = password
		}
		f.Command = "deactivate --password " + *passwordPtr
		if *forcePtr {
			f.Command = f.Command + " -f"
		}
	}
	return true
}
func (f *Flags) handleAMTInfo(amtInfoCommand *flag.FlagSet) {
	amtInfoJSONPtr := amtInfoCommand.Bool("json", false, "JSON Output")
	amtInfoMessage := AMTInfoMessage{} 
	outputMessage := ""

	amtInfoVerPtr := amtInfoCommand.Bool("ver", false, "BIOS Version")
	amtInfoBldPtr := amtInfoCommand.Bool("bld", false, "Build Number")
	amtInfoSkuPtr := amtInfoCommand.Bool("sku", false, "Product SKU")
	amtInfoUUIDPtr := amtInfoCommand.Bool("uuid", false, "Unique Identifier")
	amtInfoModePtr := amtInfoCommand.Bool("mode", false, "Current Control Mode")
	amtInfoDNSPtr := amtInfoCommand.Bool("dns", false, "Domain Name Suffix")
	amtInfoCertPtr := amtInfoCommand.Bool("cert", false, "Certificate Hashes")
	amtInfoRasPtr := amtInfoCommand.Bool("ras", false, "Remote Access Status")
	amtInfoLanPtr := amtInfoCommand.Bool("lan", false, "LAN Settings")
	amtInfoHostnamePtr := amtInfoCommand.Bool("hostname", false, "OS Hostname")
	amtInfoCommand.Parse(f.commandLineArgs[2:])

	if len(f.commandLineArgs) == 2 || (len(f.commandLineArgs) == 3 && *amtInfoJSONPtr) {
		*amtInfoVerPtr = true
		*amtInfoBldPtr = true
		*amtInfoSkuPtr = true
		*amtInfoUUIDPtr = true
		*amtInfoModePtr = true
		*amtInfoDNSPtr = true
		*amtInfoCertPtr = false
		*amtInfoRasPtr = true
		*amtInfoLanPtr = true
		*amtInfoHostnamePtr = true
	}

	if amtInfoCommand.Parsed() {
		amt := amt.Command{}
		if *amtInfoVerPtr {
			result, _ := amt.GetVersionDataFromME("AMT")
			amtInfoMessage.Version = &result
			outputMessage += "Version			: " + result + "\n"
		}
		if *amtInfoBldPtr {
			result, _ := amt.GetVersionDataFromME("Build Number")
			amtInfoMessage.BuildNumber = &result
			outputMessage += "Build Number		: " + result + "\n"
		}
		if *amtInfoSkuPtr {
			result, _ := amt.GetVersionDataFromME("Sku")
			amtInfoMessage.SKU = &result
			outputMessage += "SKU			: " + result + "\n"
		}
		if *amtInfoUUIDPtr {
			result, _ := amt.GetUUID()
			amtInfoMessage.UUID = &result
			outputMessage += "UUID			: " + result + "\n"
		}
		if *amtInfoModePtr {
			result, _ := amt.GetControlMode()
			temp := utils.InterpretControlMode(result)
			amtInfoMessage.ControlMode = &temp
			outputMessage += "Control Mode		: " + temp + "\n"
		}
		if *amtInfoDNSPtr {
			result, _ := amt.GetDNSSuffix()
			amtInfoMessage.DNSSuffix = &result
			outputMessage += "DNS Suffix		: " + result + "\n"
			result2, _ := amt.GetOSDNSSuffix()
			amtInfoMessage.DNSSuffixOS = &result2
			outputMessage += "DNS Suffix (OS)		: " + result2 + "\n"
		}
		if *amtInfoHostnamePtr {
			result, _ := os.Hostname()
			amtInfoMessage.Hostname = &result
			outputMessage += "Hostname (OS)		: " + result + "\n"
		}
		if *amtInfoRasPtr {
			result, _ := amt.GetRemoteAccessConnectionStatus()
			amtInfoMessage.RAS = &RASMessage {
				Network: result.NetworkStatus,
				RemoteStatus: result.RemoteStatus,
				Trigger: result.RemoteTrigger,
				MPSHostname: result.MPSHostname,
			}
			outputMessage += (
				"RAS Network      	: " + result.NetworkStatus + "\n" +
				"RAS Remote Status	: " + result.RemoteStatus + "\n" +
				"RAS Trigger      	: " + result.RemoteTrigger + "\n" +
				"RAS MPS Hostname 	: " + result.MPSHostname + "\n")
		}
		if *amtInfoLanPtr {
			resultWired, _ := amt.GetLANInterfaceSettings(false)
			amtInfoMessage.Wired = &LANInterfaceMessage {
				DHCPEnabled: strconv.FormatBool(resultWired.DHCPEnabled),
				DHCPMode: resultWired.DHCPMode,
				LinkStatus: resultWired.LinkStatus,
				IPAddress: resultWired.IPAddress,
				MACAddress: resultWired.MACAddress,
			}
			outputMessage += (
				"---Wired Adapter---" + "\n" +
				"DHCP Enabled 		: " + strconv.FormatBool(resultWired.DHCPEnabled) + "\n" +
				"DHCP Mode    		: " + resultWired.DHCPMode + "\n" +
				"Link Status  		: " + resultWired.LinkStatus + "\n" +
				"IP Address   		: " + resultWired.IPAddress + "\n" +
				"MAC Address  		: " + resultWired.MACAddress + "\n")

			resultWireless, _ := amt.GetLANInterfaceSettings(true)
			amtInfoMessage.Wireless = &LANInterfaceMessage {
				DHCPEnabled: strconv.FormatBool(resultWireless.DHCPEnabled),
				DHCPMode: resultWireless.DHCPMode,
				LinkStatus: resultWireless.LinkStatus,
				IPAddress: resultWireless.IPAddress,
				MACAddress: resultWireless.MACAddress,
			}
			outputMessage += (
				"---Wireless Adapter---" + "\n" +
				"DHCP Enabled 		: " + strconv.FormatBool(resultWireless.DHCPEnabled) + "\n" +
				"DHCP Mode    		: " + resultWireless.DHCPMode + "\n" +
				"Link Status  		: " + resultWireless.LinkStatus + "\n" +
				"IP Address   		: " + resultWireless.IPAddress + "\n" +
				"MAC Address  		: " + resultWireless.MACAddress + "\n")
		}
		if *amtInfoCertPtr {
			result, _ := amt.GetCertificateHashes()
			outputMessage += "Certificate Hashes	:" + "\n"
			var hashes []CertificateHashMessage
			for _, v := range result {
				hashes = append(hashes, CertificateHashMessage{
					IsDefault: v.IsDefault,
					IsActive: v.IsActive,
					Algorithm: v.Algorithm,
					Hash: v.Hash,
				})

				outputMessage += v.Name + " ("
				if v.IsDefault {
					outputMessage += "Default,"
				}
				if v.IsActive {
					outputMessage += "Active)"
				}
				outputMessage += "\n   " + v.Algorithm + ": " + v.Hash + "\n"
			}
			amtInfoMessage.CertificateHashes = &hashes
		}
		if *amtInfoJSONPtr {
			data, _ := json.MarshalIndent(amtInfoMessage, "", "\t")
			outputMessage = string(data) + "\n"
		}
		print(outputMessage)
	}
}