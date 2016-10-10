package main

import "testing"

func Test_RunBackups(t *testing.T) {
	createMockData()
	errs := RunBackups(NewFs)
	if len(errs) != 0 {
		t.Errorf("Expected '0' errors, got '%d'", len(errs))
	}
}
