package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	hostsPath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
)

var (
	// 分发名称:映射的hostname
	wslName = map[string]string{
		"Ubuntu":       "ubuntu.local",
		"Ubuntu-22.04": "ubuntu2204.local",
	}
)

func main() {

	data, err := os.ReadFile(hostsPath)
	if err != nil {
		log.Panicln(err)
	}
	var tmp bytes.Buffer
	var hostname = make(map[string][]byte, len(wslName))
	for n, h := range wslName {
		ip, err := exec.Command("wsl", "-d", n, "hostname", "-I").Output()
		if err != nil {
			log.Panicln(err)
		}
		hostname[h] = ip
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		src := strings.Split(line, " ")
		flag := false
		for _, v := range src[1:] {
			if _, ok := hostname[v]; ok {
				flag = true
				continue
			}
		}
		if flag {
			continue
		}
		tmp.WriteString(line)
		tmp.WriteString("\n")
	}

	for hostname, ip := range hostname {
		tmp.Write(bytes.TrimSpace(ip))
		tmp.WriteString(" ")
		tmp.WriteString(hostname)
		tmp.WriteString("\n")
	}

	os.WriteFile(hostsPath, tmp.Bytes(), 644)
}
