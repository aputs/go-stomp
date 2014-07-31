package stomp

import "testing"

func TestSend(t *testing.T) {
	conn, e := NewConnection("localhost", "61613")
	if e != nil {
		t.Fatal(e)
	}
	setlogger(conn)
	if e := conn.Connect(); e != nil {
		t.Fatal(e)
	}
	if e := conn.Send("/queue/test.send.01", "test data"); e != nil {
		t.Fatal(e)
	}
	if e := conn.Disconnect(); e != nil {
		t.Fatal(e)
	}
}

func TestSendWithReceipt(t *testing.T) {
	conn, e := NewConnection("localhost", "61613")
	if e != nil {
		t.Fatal(e)
	}
	setlogger(conn)
	if e := conn.Connect(); e != nil {
		t.Fatal(e)
	}
	if e := conn.Send("/queue/test.send.01", "test data", "receipt", Uuid()); e != nil {
		t.Fatal(e)
	}
	if e := conn.Disconnect(); e != nil {
		t.Fatal(e)
	}
}
