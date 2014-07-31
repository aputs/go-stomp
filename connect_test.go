package stomp

import "testing"

func TestConnectDisconnect(t *testing.T) {
	conn, e := NewConnection("localhost", "61613")
	if e != nil {
		t.Fatal(e)
	}
	setlogger(conn)
	if e := conn.Connect(); e != nil {
		t.Fatal(e)
	}
	if e := conn.Disconnect(); e != nil {
		t.Fatal(e)
	}
}
