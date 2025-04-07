package main

import "testing"

func TestFooer(t *testing.T) {
	result := Fooer(3)
	if result != "Foo" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Foo")
	}
}

func TestFooer2(t *testing.T) {
	result := Fooer2(3)
	if result != "Foo2" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Foo")
	}
}
