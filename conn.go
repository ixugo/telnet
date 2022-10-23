package telnet

import (
	"io"
	"net"
)

const (
	cr = byte('\r')
	lf = byte('\n')
)

/*
--------------------------------
IAC | Command Code | Option Code
--------------------------------
*/
const (
	CMDSE               = 240 // End of subnegotiation parameters.
	CMDNOP              = 241 // No operation.
	CMDDataMark         = 242 // The data stream portion of a Synch.
	CMDBreak            = 243 // NVT character BRK.
	CMDInterruptProcess = 244 // The function IP.
	CMDAbortOutput      = 245 // The function AO.
	CMDAreYouThere      = 246 // The function AYT.
	CMDEraseCharacter   = 247 // The function EC.
	CMDEraseLine        = 248 // The function EL.
	CMDGoAhead          = 249 // The GA signal.
	CMDSB               = 250 // Indicates that what follows is subnegotiation of the indicated option.

	CMDWILL = 251 // Accepting a request to enable.
	CMDWONT = 252 // Rejecting a request to enable.
	CMDDO   = 253 // Approving a request to enable.
	CMDDONT = 254 // Disapproving a request to enable.
	IAC     = 255 // Interpret As Command
)

const (
	optEcho            = 1 // 数据回显
	optSuppressGoAhead = 3
)

// telnetConn ...
type telnetConn struct {
	net.Conn
}

func (c *telnetConn) Close() error {
	return c.Conn.Close()
}

func (c *telnetConn) handshake(event map[[2]byte]byte) error {
	var buf [3]byte
	i, j := 0, 1
	for {
		_, err := io.ReadFull(c.Conn, buf[i:j])
		if err != nil {
			return err
		}
		// Finish
		if buf[0] == lf {
			return nil
		}
		// IAC
		if buf[i] == IAC {
			i++
			j++
			continue
		}
		// CMD
		if buf[i] == CMDWILL || buf[i] == CMDDO || buf[i] == CMDDONT || buf[i] == CMDWONT {
			i++
			j++
			continue
		}
		cmd := buf[1]
		opt := buf[2]
		// fmt.Printf("req => cmd:%d \topt:%d\n", cmd, opt)
		// OPT
		cmd, ok := event[[2]byte{cmd, opt}]
		if !ok {
			cmd = agree(buf[1])
		}
		// fmt.Printf("res => cmd:%d \topt:%d\n\n", cmd, opt)
		if err := sendCMD(c, cmd, opt); err != nil {
			return err
		}
		i, j = 0, 1
	}
}

func agree(cmd byte) byte {
	switch cmd {
	case CMDDO:
		return CMDWILL
	case CMDWILL:
		return CMDDO
	case CMDWONT:
		return CMDDONT
	case CMDDONT:
		return CMDWONT
	}
	return cmd
}

func refuse(cmd byte) byte {
	switch cmd {
	case CMDDO:
		return CMDWONT
	case CMDWILL:
		return CMDDONT
	}
	return cmd
}

func sendCMD(w io.Writer, cmd byte, opt byte) error {
	_, err := w.Write([]byte{IAC, cmd, opt})
	return err
}

func newEvent(echo bool) map[[2]byte]byte {
	data := make(map[[2]byte]byte)
	if echo {
		return data
	}
	for _, v := range [2]byte{CMDWILL, CMDDO} {
		key := [2]byte{v, optEcho}
		data[key] = refuse(v)
	}
	return data
}
