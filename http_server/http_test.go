package main

import "testing"

func TestFakeTester(t *testing.T) {
	result := FakeTester(3)
	if result != "Foo" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Foo")
	}
}
