package main

import "testing"

func TestVersionNonEmpty(t *testing.T) {
	if Version == "" {
		t.Fatal("empty Version string")
	}
}
