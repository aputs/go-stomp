package stomp_test

import (
	"testing"

	"github.com/aputs/go-stomp"
)

func TestSubscribe(t *testing.T) {
	var (
		conn *stomp.Connection
		m    <-chan stomp.Frame
		e    error
	)

	conn, e = stomp.NewConnection("localhost", "61613")
	ok(t, e)

	setlogger(conn)

	e = conn.Connect()
	ok(t, e)

	var dest = "/queue/test.subscribe"
	m, e = conn.Subscribe(dest)
	ok(t, e)

	var body = "test message"
	e = conn.Send(dest, body)
	ok(t, e)

	if e == nil {
		f := <-m
		equals(t, f.Body(), body)
	}

	e = conn.Disconnect()
	ok(t, e)
}
