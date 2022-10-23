package main

import (
	"os"
	"time"

	"github.com/ixugo/telnet"
)

func main() {
	s, err := telnet.Dial("tcp", "192.168.1.2:23", telnet.Config{Timeout: 3 * time.Second})
	if err != nil {
		panic(err)
	}
	s.Stdin = os.Stdin
	s.Stdout = os.Stdout
	if err := s.Shell(); err != nil {
		panic(err)
	}
	if err := s.Wait(); err != nil {
		panic(err)
	}
}
