package service

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"notion-native/model/entity"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

func NativeMessageService() {
	reader := bufio.NewReader(os.Stdin)

	// Chrome Native Messaging：前 4 字节是消息长度
	lengthBytes := make([]byte, 4)
	_, err := reader.Read(lengthBytes)
	if err != nil {
		return
	}

	length := binary.LittleEndian.Uint32(lengthBytes)
	messageBytes := make([]byte, length)
	reader.Read(messageBytes)

	var msg entity.Message
	json.Unmarshal(messageBytes, &msg)

	switch msg.IDE {
	case "idea":
		openWithConfig("Software\\CurrentNativeHost\\IDE", "IDEA", msg.Path)
	case "webstorm":
		openWithConfig("Software\\CurrentNativeHost\\IDE", "WEBSTORM", msg.Path)
	case "goland":
		openWithConfig("Software\\CurrentNativeHost\\IDE", "GOLAND", msg.Path)
	}

	writeNativeResponse(map[string]string{"status": "finish"})

}

// 读取Windows注册表
func readRegistryString(key, name string) (string, error) {
	var h syscall.Handle
	if err := syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER,
		syscall.StringToUTF16Ptr(key), 0, syscall.KEY_READ, &h); err != nil {
		return "", err
	}
	defer syscall.RegCloseKey(h)

	var typ uint32
	var data [1024]uint16
	dataLen := uint32(len(data) * 2) // bytes

	if err := syscall.RegQueryValueEx(h, syscall.StringToUTF16Ptr(name),
		nil, &typ, (*byte)(unsafe.Pointer(&data[0])), &dataLen); err != nil {
		return "", err
	}

	return syscall.UTF16ToString(data[:]), nil
}

func openWithConfig(regKey, ideName, filePath string) {
	idePath, err := readRegistryString(regKey, ideName)
	if err != nil || idePath == "" {
		fmt.Fprintf(os.Stderr, "IDE not configured: %s\n", ideName)
		return
	}

	if _, err := os.Stat(idePath); err != nil {
		fmt.Fprintf(os.Stderr, "IDE not found: %s\n", idePath)
		return
	}

	cmd := exec.Command(idePath, filePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	_ = cmd.Start()
}

func writeNativeResponse(resp any) {
	b, _ := json.Marshal(resp)
	length := uint32(len(b))
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, length)
	os.Stdout.Write(buf)
	os.Stdout.Write(b)
}
