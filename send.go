package stomp

import (
	"errors"
	"net"
)

func (c *Connection) Send(dest string, body string, headers ...string) error {
	c.log("start send...")

	f := NewFrame(SEND)
	f.AddHeader("destination", dest)
	for k, v := range makeheaders(headers) {
		f.AddHeader(k, v)
	}

	var receiptRequested bool

	if _, found := f.Headers()["receipt"]; found {
		receiptRequested = true
	}

	rf, e := c.writeFrame(f)
	if e != nil {
		if neterr, ok := e.(net.Error); ok && neterr.Timeout() && receiptRequested {
			return e
		}
	}

	if receiptRequested && rf.command != RECEIPT {
		return errors.New("receipt was requested, but was not received")
	}

	c.log("done sending...")
	return nil
}
