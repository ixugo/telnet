package main

import (
	"time"

	"github.com/ixugo/telnet"
)

func main() {
	s, err := telnet.Dial("tcp", "192.168.1.2:23", telnet.Config{Timeout: 3 * time.Second})
	if err != nil {
		panic(err)
	}
	// 自定义输入
	w, err := s.StdinPipe()
	if err != nil {
		panic(err)
	}
	_ = w
	// 自定义输出
	r, err := s.StdoutPipe()
	if err != nil {
		panic(err)
	}
	_ = r

	if err := s.Shell(); err != nil {
		panic(err)
	}

	if err := s.Wait(); err != nil {
		panic(err)
	}
}
