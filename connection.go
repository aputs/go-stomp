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

	in            chan *Frame
	out           chan *Frame
	err           chan ConnectionError
	subscriptions map[string]chan Frame

	connected         bool
	connectionClosing chan bool

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

	c.in = make(chan *Frame)
	c.out = make(chan *Frame)
	c.err = make(chan ConnectionError)
	c.subscriptions = make(map[string]chan Frame)

	c.connectionClosing = make(chan bool, 1)
	c.connected = true

	go c.connectionHandler()
	go c.errorHandler()
	go c.outgoing()
	go c.incoming()

	return
}

func (c *Connection) connectionHandler() {
	for {
		select {
		case <-c.connectionClosing:
			c.connected = false
		}
	}
}

func (c *Connection) errorHandler() {
	for c.connected {
		select {
		case e := <-c.err:
			c.log(fmt.Sprintf("connection error: %s", e.e))
		}
	}
}

func (c *Connection) incoming() {
	for c.connected {
		buf := make([]byte, 4096)
		n, e := c.Conn.Read(buf)
		if e != nil {
			if neterr, ok := e.(net.Error); ok && !neterr.Timeout() {
				c.err <- ConnectionError{c: c, e: e, f: NullFrame}
			}
			continue
		}

		f, e := ParseFrame(buf)
		if e != nil {
			c.err <- ConnectionError{c: c, e: e, f: NullFrame}
			continue
		}

		c.log(fmt.Sprintf("received %q", string(buf[:n])))

		switch f.command {
		case MESSAGE:
			dest, _ := f.GetHeader("destination")
			if _, found := c.subscriptions[dest]; found {
				c.log(fmt.Sprintf("% #v", f))
				c.subscriptions[dest] <- *f
				break
			}
			fallthrough
		default:
			c.in <- f
		}

	}

	c.Close()
}

func (c *Connection) outgoing() {
	for c.connected {
		select {
		case out := <-c.out:
			c.log(fmt.Sprintf("sending %q", string(out.Bytes())))
			_, e := c.Conn.Write(out.Bytes())
			if e != nil {
				c.err <- ConnectionError{e: e, c: c, f: out}
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

func (c *Connection) Connected() bool {
	return c.connected
}

func (c *Connection) SetLogger(l *log.Logger) {
	c.logger = l

}
