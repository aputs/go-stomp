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
	c.ResponseTimeout = 10 * time.Second

	// TODO global logger
	c.log(fmt.Sprintf("tcp connection opened %s", c.Conn.RemoteAddr().String()))
	return
}

func (c *Connection) writeFrame(f Frame) (*Frame, error) {
	c.Lock()
	defer c.Unlock()

	c.Conn.Write(f.Bytes())
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

func (c *Connection) Connect(headers ...string) error {
	c.log("connecting...")

	f := NewFrame(CONNECT)
	for k, v := range makeheaders(headers) {
		f.AddHeader(k, v)
	}

	// sanity checks
	if _, found := f.headers["accept-version"]; !found {
		f.AddHeader("accept-version", "1.1")
	}
	if _, found := f.headers["host"]; !found {
		f.AddHeader("host", c.Conn.LocalAddr().String())
	}

	// TODO setup heartbeating
	rf, e := c.writeFrame(f)
	if e != nil {
		return e
	}
	if rf.command != CONNECTED {
		return errors.New("connection error")
	}
	if rf.command == ERROR {
		return errors.New(rf.Body())
	}

	c.session = rf.headers["session"]
	c.version = rf.headers["version"]
	c.server = rf.headers["server"]
	c.log("connected.")
	return nil
}

func (c *Connection) Disconnect() error {
	c.log("disconnecting...")
	f := NewFrame(DISCONNECT)
	if _, e := c.writeFrame(f); e != nil {
		c.log(e)
		return e
	}
	c.log("disconnected.")
	return nil
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
