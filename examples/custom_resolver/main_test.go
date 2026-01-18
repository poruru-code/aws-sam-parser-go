// Where: examples/custom_resolver/main_test.go
// What: Smoke test for the custom resolver example.
// Why: Ensure the example stays runnable.
package main

import "testing"

func TestRun(t *testing.T) {
	if err := run(); err != nil {
		t.Fatalf("run() error: %v", err)
	}
}
