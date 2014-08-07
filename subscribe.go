package stomp

import (
	"errors"
	"time"
)

func (c *Connection) Subscribe(dest string, headers ...string) (<-chan Frame, error) {
	c.log("start subscribe...")

	f := NewFrame(SUBSCRIBE)
	f.AddHeader("destination", dest)

	if c, found := c.subscriptions[dest]; found {
		return c, nil
	}

	for k, v := range makeheaders(headers) {
		f.AddHeader(k, v)
	}

	if _, found := f.GetHeader("id"); found {
		f.AddHeader("id", Uuid())
	}

	if x, found := f.GetHeader("ack"); found {
		if _, valid := map[string]bool{
			"auto":              true,
			"client":            true,
			"client-individual": true,
		}[x]; !valid {
			return nil, errors.New("unknown ack value")
		}
	} else {
		f.AddHeader("ack", "auto")
	}

	var (
		receiptRequested bool
		receiptReceived  bool
	)

	if _, found := f.GetHeader("receipt"); found {
		receiptRequested = true
	}

	c.out <- f

	select {
	case rf := <-c.in:
		switch rf.command {
		case RECEIPT:
			receiptReceived = true
		case ERROR:
			msg, _ := rf.GetHeader("message")
			return nil, errors.New(msg)
		}
	case <-time.After(c.ResponseTimeout):
		break
	}

	if receiptRequested && !receiptReceived {
		return nil, errors.New("receipt was requested, but was not received")
	}

	rc := make(chan Frame)

	c.subscriptions[dest] = rc

	c.log("done subscribing...")

	return rc, nil
}
