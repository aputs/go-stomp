package stomp

import (
	"errors"
	"strings"
)

const (
	NULL byte = 0
	LF   byte = 13
	CR   byte = 10
	EOL       = CR
)

type Headers map[string]string

func makeheaders(hs []string) map[string]string {
	h := map[string]string{}
	// pad to even size
	if (len(hs) % 2) != 0 {
		hs = append(hs, "")
	}
	for i := 0; i < len(hs); i += 2 {
		// ignore keys already in the hash
		if _, found := h[hs[i]]; found {
			continue
		}
		h[hs[i]] = hs[i+1]
	}
	return h
}

type Command int

const (
	UNKNOWN Command = iota
	CONNECTED
	MESSAGE
	RECEIPT
	ERROR
	SEND
	SUBSCRIBE
	UNSUBSCRIBE
	BEGIN
	COMMIT
	ABORT
	ACK
	NACK
	DISCONNECT
	CONNECT
	STOMP
)

var commands = map[string]Command{
	"CONNECTED":   CONNECTED,
	"MESSAGE":     MESSAGE,
	"RECEIPT":     RECEIPT,
	"ERROR":       ERROR,
	"SEND":        SEND,
	"SUBSCRIBE":   SUBSCRIBE,
	"UNSUBSCRIBE": UNSUBSCRIBE,
	"BEGIN":       BEGIN,
	"COMMIT":      COMMIT,
	"ABORT":       ABORT,
	"ACK":         ACK,
	"NACK":        NACK,
	"DISCONNECT":  DISCONNECT,
	"CONNECT":     CONNECT,
	"STOMP":       STOMP,
}

func (c Command) String() string {
	if s, found := map[Command]string{
		CONNECTED:   "CONNECTED",
		MESSAGE:     "MESSAGE",
		RECEIPT:     "RECEIPT",
		ERROR:       "ERROR",
		SEND:        "SEND",
		SUBSCRIBE:   "SUBSCRIBE",
		UNSUBSCRIBE: "UNSUBSCRIBE",
		BEGIN:       "BEGIN",
		COMMIT:      "COMMIT",
		ABORT:       "ABORT",
		ACK:         "ACK",
		NACK:        "NACK",
		DISCONNECT:  "DISCONNECT",
		CONNECT:     "CONNECT",
		STOMP:       "STOMP",
	}[c]; found {
		return s
	}
	return "UNKNOWN"
}

type Frame struct {
	command Command
	headers map[string]string
	body    string
}

var NullFrame = &Frame{command: UNKNOWN, headers: map[string]string{}}

func NewFrame(cmd Command) Frame {
	return Frame{command: cmd}
}

func (f *Frame) AddBody(body string) *Frame {
	f.body = f.body + body
	return f
}

func (f *Frame) Body() string {
	return f.body
}

func (f *Frame) AddHeader(key, value string) *Frame {
	if f.headers == nil {
		f.headers = Headers{}
	}
	// TODO check header value sanity
	f.headers[key] = value
	return f
}

func (f *Frame) Headers() map[string]string {
	return f.headers
}

func (f *Frame) Bytes() []byte {
	var buf []byte

	eol := []byte{EOL}
	cmd := []byte(f.command.String())
	buf = append(buf, cmd...)
	buf = append(buf, eol...)

	for k, v := range f.headers {
		buf = append(buf, k...)
		buf = append(buf, ':')
		buf = append(buf, v...)
		buf = append(buf, eol...)
	}

	buf = append(buf, eol...)
	buf = append(buf, NULL)
	buf = append(buf, eol...)
	return buf
}

func ParseFrame(b []byte) (*Frame, error) {
	p := []string{}
	buf := []byte{}

	for i := 0; i < len(b); i++ {
		switch b[i] {
		case CR:
			p = append(p, string(buf))
			buf = []byte{}
		case NULL:
			p = append(p, "")
		default:
			buf = append(buf, b[i])
		}
	}

	if len(p) < 4 {
		return nil, errors.New("invalid frame!")
	}

	c, found := commands[p[0]]
	if !found {
		return nil, errors.New("invalid frame!")
	}

	of := NewFrame(c)
	of.AddBody(p[len(p)-1])
	for _, h := range p[1 : len(p)-2] {
		if kv := strings.Split(h, ":"); len(kv) == 2 {
			of.AddHeader(kv[0], kv[1])
		}
	}

	return &of, nil
}
