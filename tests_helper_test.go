package stomp

import (
	"log"
	"os"
	"testing"
)

type DevNull struct{}

func (DevNull) Write(p []byte) (int, error) {
	return len(p), nil
}

func setlogger(conn *Connection) {
	if testing.Verbose() {
		conn.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds))
	} else {
		conn.SetLogger(log.New(new(DevNull), "", 0))
	}
}
