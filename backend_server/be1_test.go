package main

import (
	"testing"
	"time"
)

func TestFooer(t *testing.T) {
	result := Fooer(3)
	time.Sleep(10 * time.Second)
	if result != "Foo" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Foo")
	}
}

func TestFooer2(t *testing.T) {
	result := Fooer2(6)
	time.Sleep(10 * time.Second)
	if result != "Foo" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Foo")
	}
}
