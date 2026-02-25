package service

import (
	"log"
	"net/url"
	"os"
	"os/exec"
)

func NativeProtocolService() {

	if len(os.Args) < 2 {
		return
	}

	raw := os.Args[1]
	log.Println("raw url:", raw)

	u, err := url.Parse(raw)
	if err != nil {
		log.Println("parse error:", err)
		return
	}

	action := u.Host // open
	q := u.Query()

	typ := q.Get("type")
	path := q.Get("path")

	if path == "" {
		log.Println("empty path")
		return
	}

	// URL decode
	path, _ = url.QueryUnescape(path)

	log.Println("action:", action)
	log.Println("type:", typ)
	log.Println("path:", path)

	switch typ {
	case "file", "folder":
		openWithExplorer(path)

	default:
		log.Println("unknown type:", typ)
	}
}

func openWithExplorer(path string) {
	// explorer 对文件 & 文件夹都有效
	exec.Command("explorer", path).Start()
}
