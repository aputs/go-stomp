package stomp

import "errors"

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
	buf = append(buf, f.body...)
	buf = append(buf, NULL)
	//buf = append(buf, eol...)

	return buf
}

func ParseFrame(buf []byte) (*Frame, error) {
	if len(buf) == 0 {
		return nil, errors.New("invalid frame!")
	}

	of := &Frame{headers: map[string]string{}}
	eol := byte(EOL)
	pos := 0
	epos := pos

	// frame type
	for ; buf[epos] != NULL; epos++ {
		if buf[epos] == eol {
			break
		}
	}
	ft := string(buf[pos:epos])
	found := false
	if of.command, found = commands[ft]; !found {
		return nil, errors.New("unknown frame type!")
	}

	// headers
	epos++
	for buf[epos] != NULL {
		pos = epos
		cpos := pos
		for ; buf[epos] != NULL; epos++ {
			if buf[epos] == eol {
				break
			}
			if buf[epos] == ':' {
				cpos = epos
			}
		}
		hk := string(buf[pos:cpos])
		if len(hk) == 0 {
			break
		}
		hv := string(buf[cpos+1 : epos])
		of.AddHeader(hk, hv)
		epos++
	}

	// body
	epos++
	pos = epos
	for ; buf[epos] != NULL; epos++ {
	}
	of.AddBody(string(buf[pos:epos]))

	return of, nil
}
