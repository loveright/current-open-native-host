package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"net/url"
	"notion-native/service"
	"os"
	"strings"
)

func main() {

	if isProtocolLaunch() {
		service.NativeProtocolService()
		return
	}

	if len(os.Args) >= 3 && os.Args[2] == "generate" {

		path := os.Args[1]

		mode := "default"
		if len(os.Args) >= 4 {
			mode = os.Args[3]
		}

		generateLocationLink(path, mode)
		return
	}

	service.NativeMessageService()
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
