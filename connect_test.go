package stomp

import (
	"log"
	"os"
	"testing"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

func TestConnectDisconnect(t *testing.T) {
	conn, e := NewConnection("localhost", "61613")
	if e != nil {
		t.Fatal(e)
	}
	conn.SetLogger(logger)
	conn.Connect("user", "guest")
	conn.Disconnect()
}
