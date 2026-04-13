package service

import (
	"encoding/binary"
	"encoding/json"
	"notion-native/model/entity"
	"os"
	"os/exec"
)

func NativeService() {
	// 1. 读长度
	lengthBytes := make([]byte, 4)
	if _, err := os.Stdin.Read(lengthBytes); err != nil {
		return
	}
	length := binary.LittleEndian.Uint32(lengthBytes)

	// 2. 读消息
	messageBytes := make([]byte, length)
	if _, err := os.Stdin.Read(messageBytes); err != nil {
		return
	}

	var msg entity.Message
	json.Unmarshal(messageBytes, &msg)

	if msg.IDE == "idea" {
		openIdea(msg.Path)
	}

	// 3. 必须回一条消息给 Chrome
	resp := map[string]string{"status": "ok"}
	respBytes, _ := json.Marshal(resp)

	// 写长度
	binary.Write(os.Stdout, binary.LittleEndian, uint32(len(respBytes)))
	// 写内容
	os.Stdout.Write(respBytes)

	// 4. 立刻退出（非常重要）

}

func openIdea(path string) error {
	cmd := exec.Command(
		"cmd",
		"/c",
		`D:\E\software\software1\IDEA\install\IntelliJ IDEA 2025.1\bin\idea.bat`,
		path,
	)
	return cmd.Start()
}
