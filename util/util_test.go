package util

import "testing"

func TestDebug(t *testing.T) {
	Verbose = true
	Debug("oi")
}
