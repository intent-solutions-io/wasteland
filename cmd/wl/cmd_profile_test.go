package main

import (
	"bytes"
	"testing"
)

func TestPrintBar_Normal(t *testing.T) {
	var buf bytes.Buffer
	printBar(&buf, "Test", 0.5)
	out := buf.String()
	if out == "" {
		t.Fatal("expected output")
	}
	if !contains(out, "50%") {
		t.Errorf("expected 50%% in output, got: %q", out)
	}
}

func TestPrintBar_Zero(t *testing.T) {
	var buf bytes.Buffer
	printBar(&buf, "Test", 0)
	if buf.Len() == 0 {
		t.Fatal("expected output")
	}
}

func TestPrintBar_One(t *testing.T) {
	var buf bytes.Buffer
	printBar(&buf, "Test", 1.0)
	out := buf.String()
	if !contains(out, "100%") {
		t.Errorf("expected 100%% in output, got: %q", out)
	}
}

func TestPrintBar_NegativeClamps(t *testing.T) {
	var buf bytes.Buffer
	// Should not panic
	printBar(&buf, "Test", -0.5)
	out := buf.String()
	if !contains(out, "0%") {
		t.Errorf("expected 0%% in output for negative value, got: %q", out)
	}
}

func TestPrintBar_OverOneClamps(t *testing.T) {
	var buf bytes.Buffer
	// Should not panic
	printBar(&buf, "Test", 1.5)
	out := buf.String()
	if !contains(out, "100%") {
		t.Errorf("expected 100%% in output for >1 value, got: %q", out)
	}
}

func TestClamp01(t *testing.T) {
	tests := []struct {
		in, want float64
	}{
		{0.5, 0.5},
		{0, 0},
		{1, 1},
		{-1, 0},
		{2, 1},
		{-0.001, 0},
		{1.001, 1},
	}
	for _, tc := range tests {
		got := clamp01(tc.in)
		if got != tc.want {
			t.Errorf("clamp01(%f) = %f, want %f", tc.in, got, tc.want)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && bytes.Contains([]byte(s), []byte(sub))
}
