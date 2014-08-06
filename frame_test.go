package stomp

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

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
	var buf = []byte("MESSAGE\ndestination:/queue/test.subscribe\n\ntest message\x00")

	f, _ := ParseFrame(buf)
	equals(t, buf, []byte(string(f.Bytes())))
}
