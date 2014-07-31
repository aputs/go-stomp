package stomp

import (
	"errors"
	"net"
)

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

	// TODO setup heartbeating
	rf, e := c.writeFrame(f)
	if e != nil {
		return e
	}
	switch rf.command {
	case CONNECTED:
	case ERROR:
		return errors.New(rf.Headers()["message"])
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
	_, e := c.writeFrame(f)
	if e != nil {
		if neterr, ok := e.(net.Error); ok && !neterr.Timeout() {
			return e
		}
	}
	c.log("disconnected.")
	return nil
}
