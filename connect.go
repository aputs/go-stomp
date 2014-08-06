package stomp

import (
	"errors"
	"time"
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

	c.out <- f

	select {
	case rf := <-c.in:
		switch rf.command {
		case CONNECTED:
			c.session = rf.headers["session"]
			c.version = rf.headers["version"]
			c.server = rf.headers["server"]
			return nil

			c.log("connected.")
		case ERROR:
			return errors.New(rf.Headers()["message"])
		}
	case <-time.After(c.ResponseTimeout):
		break
	}

	return errors.New("No Response was received.")
}

func (c *Connection) Disconnect() error {
	c.log("disconnecting...")

	c.out <- NewFrame(DISCONNECT)

	select {
	case rf := <-c.in:
		switch rf.command {
		case CONNECTED:
			return nil

			c.log("connected.")
		case ERROR:
			return errors.New(rf.Headers()["message"])
		}
	case <-time.After(c.ResponseTimeout):
		break
	}

	c.Close()

	c.log("disconnected.")

	c.session = ""
	c.version = ""
	c.server = ""

	return nil
}
