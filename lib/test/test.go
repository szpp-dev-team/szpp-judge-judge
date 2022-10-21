package test

import (
	"testing"

	"golang.org/x/exp/constraints"
)

func Assert(t *testing.T, b bool) {
	if !b {
		t.Fatal()
	}
}

// check if left == right
func AssertEq[T constraints.Ordered](t *testing.T, left T, right T) {
	if left != right {
		t.Fatalf("%v != %v", left, right)
	}
}

// check if left != right
func AssertNeq[T constraints.Ordered](t *testing.T, left T, right T) {
	if left == right {
		t.Fatalf("%v == %v", left, right)
	}
}

// check if left > right
func AssertGt[T constraints.Ordered](t *testing.T, left T, right T) {
	if left <= right {
		t.Fatalf("%v <= %v", left, right)
	}
}
