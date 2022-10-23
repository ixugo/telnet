package telnet

import (
	"bytes"
	"io"
)

// Websocket ...
type Websocket interface {
	WriteMessage(messageType int, data []byte)
}

// NewStdout 提供简单低效的转换
func NewStdout(websocket Websocket) io.Writer {
	buf := new(bytes.Buffer)
	go func() {
		for {
			r, size, err := buf.ReadRune()
			if err != nil {
				buf.Reset()
				return
			}
			if size > 0 {
				websocket.WriteMessage(1, []byte(string(r)))
			}
		}
	}()
	return buf
}
