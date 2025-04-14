package main

import (
	"testing"
	"time"
)

func TestFakeTester(t *testing.T) {
	result := FakeTester(3)
	time.Sleep(15 * time.Second)

	if result != "Foo" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Foo")
	}
}
