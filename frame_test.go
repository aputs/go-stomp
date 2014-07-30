package stomp

import "testing"

func TestMakeheaders(t *testing.T) {
	// invalid sized headers, must not panic
	makeheaders([]string{"test", "124", "te"})
}
