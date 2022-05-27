package test

import "testing"

func Assert(t *testing.T, expected interface{}, result interface{}) {
	if expected != result {
		t.Fatalf("Expected: %v but instead got %v", expected, result)
	}
}
