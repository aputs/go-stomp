package stomp

import (
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
	in      chan Frame
	out     chan Frame
	err     chan ConnectionError

	logger *log.Logger
}

type ConnectionError struct {
	c *Connection
	f *Frame
	e error
}

func NewConnection(host, port string) (c *Connection, e error) {
	c = &Connection{}
	c.Conn, e = net.Dial("tcp", net.JoinHostPort(host, port))
	if e != nil {
		return nil, e
	}

	c.session = "none"
	c.ResponseTimeout = 200 * time.Millisecond

	c.in = make(chan Frame)
	c.out = make(chan Frame)
	c.err = make(chan ConnectionError)

	go c.outgoing()
	go c.incoming()

	return
}

func (c *Connection) incoming() {
	for {
		buf := make([]byte, 4096)
		n, e := c.Conn.Read(buf)
		if e != nil {
			if neterr, ok := e.(net.Error); ok && neterr.Timeout() {
				c.err <- ConnectionError{c: c, e: e, f: NullFrame}
				continue
			}
		}

		f, e := ParseFrame(buf)
		if e != nil {
			c.err <- ConnectionError{c: c, e: e, f: NullFrame}
			continue
		}

		c.log(fmt.Sprintf("received %q", string(buf[:n])))
		c.in <- *f
	}
}

func (c *Connection) outgoing() {
	for {
		select {
		case out := <-c.out:
			c.log(fmt.Sprintf("sending %q", string(out.Bytes())))
			_, e := c.Conn.Write(out.Bytes())
			if e != nil {
				c.err <- ConnectionError{e: e, c: c, f: &out}
			}
		}
	}
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
