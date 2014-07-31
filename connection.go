package stomp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Connection struct {
	ResponseTimeout time.Duration

	sync.Mutex
	net.Conn
	session string
	version string
	server  string
	logger  *log.Logger
}

func NewConnection(host, port string) (c *Connection, e error) {
	c = &Connection{}
	c.Conn, e = net.Dial("tcp", net.JoinHostPort(host, port))
	if e != nil {
		return nil, e
	}

	c.session = "none"
	c.ResponseTimeout = 200 * time.Millisecond

	return
}

func (c *Connection) writeFrame(f Frame) (*Frame, error) {
	c.Lock()
	defer c.Unlock()

	var n int

	var buf = f.Bytes()

	n, e := c.Conn.Write(buf)
	if e != nil {
		return nil, e
	}

	if n != len(buf) {
		return nil, errors.New("not all bytes were written")
	}

	var resp []byte

	status := make(chan error)
	go func() {
		for {
			buf := make([]byte, 4096)
			c.Conn.SetReadDeadline(time.Now().Add(c.ResponseTimeout))
			n, e := c.Conn.Read(buf)
			if e != nil {
				status <- e
				return
			}
			for i := 0; i < n; i++ {
				resp = append(resp, buf[i])
				if buf[i] == NULL {
					status <- nil
					return
				}
			}
		}
	}()

	select {
	case err := <-status:
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return nil, err
		}
	}

	if len(resp) == 0 {
		uf := NewFrame(UNKNOWN)
		return &uf, nil
	}

	return ParseFrame(resp)
}

// connection logging
func (c *Connection) log(v ...interface{}) {
	if c.logger != nil {
		c.logger.Printf("[%s] %s", c.session, fmt.Sprint(v...))
	}
	return
}

func (c *Connection) SetLogger(l *log.Logger) {
	c.logger = l

}
