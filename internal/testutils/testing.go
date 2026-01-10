/*
Package testutils provides simple utility functions for testing.
*/
package testutils

import (
	"strings"
	"testing"
)

// require wraps testing.TB to provide FailNow on Errorf.
type require struct {
	testing.TB
}

// Errorf logs the formatted error message and fails the test immediately.
func (r require) Errorf(format string, args ...any) {
	r.Logf(format, args...)
	r.FailNow()
}

// Require returns a testing.TB that fails immediately on errors.
func Require(t testing.TB) testing.TB {
	t.Helper()
	return require{TB: t}
}

// Contains asserts that str contains substr.
func Contains(t testing.TB, str, substr string) {
	t.Helper()
	if !strings.Contains(str, substr) {
		t.Errorf("expected '%s' to contain '%s'", str, substr)
	}
}

// Equal asserts that expected and actual are equal.
func Equal[T comparable](t testing.TB, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// Error asserts that err is not nil.
func Error(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// False asserts that the condition is false.
func False(t testing.TB, condition bool) { //nolint:revive // input condition to control flow
	t.Helper()
	if condition {
		t.Error("expected false, got true")
	}
}

// NoError asserts that err is nil.
func NoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error, got '%s'", err)
	}
}

// NotContains asserts that str does not contain substr.
func NotContains(t testing.TB, str, substr string) {
	t.Helper()
	if strings.Contains(str, substr) {
		t.Errorf("expected '%s' not to contain '%s'", str, substr)
	}
}

// NotNil asserts that the object is not nil.
func NotNil(t testing.TB, object any) {
	t.Helper()
	if object == nil {
		t.Error("expected not nil, got nil")
	}
}

// True asserts that the condition is true.
func True(t testing.TB, condition bool) { //nolint:revive // input condition to control flow
	t.Helper()
	if !condition {
		t.Error("expected true, got false")
	}
}
