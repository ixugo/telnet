// Package telnet 参考 https://www.rfc-editor.org/rfc/rfc854
package telnet

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

// Session 实现了 telnet 连接
type Session struct {
	conn                  telnetConn
	Stdin                 io.Reader
	Stdout                io.Writer
	started               bool
	stdinpipe, stdoutpipe bool
	errors                chan error
}

// NewSession ...
func NewSession(conn telnetConn) *Session {
	return &Session{conn: conn}
}

type Config struct {
	Timeout time.Duration
	Echo    bool
}

// Dial 与指定的 telnet 服务器连接
func Dial(network, addr string, cfg Config) (*Session, error) {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}
	conn, err := net.DialTimeout(network, addr, cfg.Timeout)
	if err != nil {
		return nil, err
	}
	c := telnetConn{Conn: conn}
	err = c.handshake(newEvent(cfg.Echo))
	return NewSession(c), err
}

// Shell ...
func (s *Session) Shell() error {
	if s.started {
		return errors.New("telnet: session already started")
	}
	s.started = true
	s.errors = make(chan error, 2)
	type F func() error
	for _, fn := range []F{s.stdin(), s.stdout()} {
		go func(fn func() error) {
			if fn == nil {
				return
			}
			s.errors <- fn()
		}(fn)
	}
	return nil
}

func (s *Session) StdoutPipe() (io.Reader, error) {
	if s.Stdout != nil {
		return nil, errors.New("telnet: Stdin already set")
	}
	if s.started {
		return nil, errors.New("telnet: StdoutPipe after process started")
	}
	s.stdoutpipe = true
	return s.conn, nil
}

func (s *Session) StdinPipe() (io.WriteCloser, error) {
	if s.Stdin != nil {
		return nil, errors.New("telnet: Stdin already set")
	}
	if s.started {
		return nil, errors.New("telnet: StdinPipe after process started")
	}
	s.stdinpipe = true
	r, w := io.Pipe()
	go func() {
		_, err := io.Copy(s.conn, r)
		w.CloseWithError(err)
	}()
	return w, nil
}

func (s *Session) stdin() func() error {
	if s.stdinpipe {
		return nil
	}
	var stdin io.Reader = s.Stdin
	if stdin == nil {
		stdin = new(bytes.Buffer)
	}
	return func() error {
		_, err := io.Copy(s.conn, stdin)
		return err
	}
}

func (s *Session) stdout() func() error {
	if s.stdoutpipe {
		return nil
	}
	if s.Stdout == nil {
		s.Stdout = io.Discard
	}
	return func() error {
		_, err := io.Copy(s.Stdout, s.conn)
		return err
	}
}

func (s *Session) Wait() error {
	if !s.started {
		return fmt.Errorf("telnet: session not started")
	}
	err := <-s.errors
	return err
}

func (s *Session) Close() error {
	close(s.errors)
	return s.conn.Close()
}
