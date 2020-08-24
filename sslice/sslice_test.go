package sslice

import (
	"testing"
)

func TestAppendAndRemoveStr(t *testing.T) {
	x := make([]string, 5)
	x[0] = "a"
	x[1] = "b"
	x[2] = "c"
	x[3] = "d"
	x[4] = "e"

	AppendStr(x, "f", "g")

	RemoveStr(x, "d")

}
