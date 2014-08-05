package stomp_test

import (
	"testing"

	"github.com/aputs/go-stomp"
)

func TestSend(t *testing.T) {
	var (
		conn *stomp.Connection
		e    error
	)

	conn, e = stomp.NewConnection("localhost", "61613")
	ok(t, e)

	setlogger(conn)

	e = conn.Connect()
	ok(t, e)

	e = conn.Send("/queue/test.send.01", "test data")
	ok(t, e)

	e = conn.Disconnect()
	ok(t, e)
}

func TestSendWithReceipt(t *testing.T) {
	var (
		conn *stomp.Connection
		e    error
	)

	conn, e = stomp.NewConnection("localhost", "61613")
	ok(t, e)

	setlogger(conn)

	e = conn.Connect()
	ok(t, e)

	e = conn.Send("/queue/test.send.01", "test data", "receipt", stomp.Uuid())
	ok(t, e)

	e = conn.Disconnect()
	ok(t, e)
}
