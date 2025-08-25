package main

import (
	"testing"
	"time"
)

func TestSSHLocalhost(t *testing.T) {
	if r, err := scanSSHPortWithErr("noregdev8", 22, time.Millisecond*100); err != nil {
		t.Errorf("%+v", err)
	} else {
		t.Logf("Result: %v", r)
	}
}
