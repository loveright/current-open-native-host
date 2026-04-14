package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gin-gonic/gin"
	"net"
	"net/url"
	"notion-native/config"
	"notion-native/routers"
	"notion-native/service"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("args:", os.Args)
	// ===== 1. 如果是 server 进程，只跑 HTTP =====
	if isServerMode() {
		startServer()
		return
	}

	// ===== 2. 执行一次性任务 =====
	runOnceTask()

	// ===== 3. 如果服务没启动，就拉起子进程 =====
	if !isPortInUse(config.Port) {
		fmt.Println("启动子进程 server...")
		startServerProcess()
	} else {
		fmt.Println("server 已存在")
	}

}

func isProtocolLaunch() bool {
	// 协议启动一定会带参数
	// native messaging 启动时是没有参数的
	return len(os.Args) > 1 &&
		strings.HasPrefix(os.Args[1], "current-opener://")
}

func generateLocationLink(path string, mode string) {

	encoded := url.QueryEscape(path)

	typ := "file"

	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		typ = "folder"
	}

	if mode == "idea" || mode == "webstorm" || mode == "goland" {
		typ = mode
	}

	link := fmt.Sprintf(
		"https://current-opener/open?type=%s&path=%s",
		typ,
		encoded,
	)

	markdown := fmt.Sprintf("[%s](%s)", path, link)

	clipboard.WriteAll(markdown)
}

func isPortInUse(port string) bool {
	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		return false
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	return true
}

func isServerMode() bool {
	return len(os.Args) > 1 && os.Args[1] == "server"
}

func startServerProcess() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取 exePath 失败:", err)
		return
	}

	fmt.Println("exePath:", exePath)

	cmd := exec.Command(exePath, "server")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil
	// Windows 可选隐藏窗口
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err = cmd.Start()
	if err != nil {
		fmt.Println("启动子进程失败:", err)
		return
	}

	fmt.Println("子进程 PID:", cmd.Process.Pid)
}

func startServer() {
	fmt.Println("HTTP server 启动中...")
	r := gin.Default()
	routers.RegisterRouters(r)

	err := r.Run(config.Port)
	if err != nil {
		fmt.Println("HTTP 启动失败:", err)
		panic(err)
	}
}

func runOnceTask() bool {

	if isProtocolLaunch() {
		service.NativeProtocolService()
		return true
	}

	if len(os.Args) >= 3 && os.Args[2] == "generate" {
		path := os.Args[1]
		mode := "default"
		if len(os.Args) >= 4 {
			mode = os.Args[3]
		}

		generateLocationLink(path, mode)
		return true
	}

	if isNativeMessageLaunch() {
		service.NativeMessageService()
		return true
	}
	return false
}

func isNativeMessageLaunch() bool {
	return len(os.Args) > 1 && os.Args[1] == "native"
}
