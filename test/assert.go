package test

import "testing"

func AssertEqual(t *testing.T, expected interface{}, result interface{}) {
	if expected != result {
		t.Fatalf("Expected: %v but instead got %v", expected, result)
	}
}

func AssertNil(t *testing.T, value interface{}) {
	if value != nil {
		t.Fatalf("Expected value to be nil but instead got: %v", value)
	}
}

func AssertDefined(t *testing.T, value interface{}) {
	if value == nil {
		t.Fatal("Expected value to be defined but instead got nil")
	}
}
