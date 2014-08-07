package stomp

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// ok fails the test if an err is not nil.
func nok(tb testing.TB, err error) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: expected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestMakeheaders(t *testing.T) {
	// invalid sized headers, must not panic
	makeheaders([]string{"test", "124", "te"})
}

func TestParseFrame_Message(t *testing.T) {
	var buf = []byte("MESSAGE\ndestination:/queue/test.subscribe\n\ntest 界\nmessage\x00")

	_, e := ParseFrame(buf)
	ok(t, e)
}

func TestValue_encode(t *testing.T) {
	equals(t, "t\\reКакая-то\\\\s\\ct\\n界", encode("t\reКакая-то\\s:t\n界"))
}

func TestValue_decode(t *testing.T) {
	var (
		d string
		e error
	)
	d, e = decode("界\\r124\\n\\c\\")
	ok(t, e)
	equals(t, "界\r124\n:\\", d)
	_, e = decode("\\x")
	nok(t, e)
	_, e = decode("xx\\|n")
	nok(t, e)
}
