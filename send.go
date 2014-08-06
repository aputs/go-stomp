package stomp

import (
	"errors"
	"time"
)

func (c *Connection) Send(dest string, body string, headers ...string) error {
	c.log("start send...")

	f := NewFrame(SEND)
	f.AddHeader("destination", dest)
	for k, v := range makeheaders(headers) {
		f.AddHeader(k, v)
	}

	var (
		receiptRequested bool
		receiptReceived  bool
	)

	if _, found := f.Headers()["receipt"]; found {
		receiptRequested = true
	}

	f.AddBody(body)

	c.out <- f

	if receiptRequested {
		select {
		case rf := <-c.in:
			switch rf.command {
			case RECEIPT:
				receiptReceived = true
			case ERROR:
				return errors.New(rf.Headers()["message"])
			}
		case <-time.After(c.ResponseTimeout):
			break
		}
		if !receiptReceived {
			return errors.New("receipt was requested, but was not received")
		}
	}

	c.log("done sending...")

	return nil
}
