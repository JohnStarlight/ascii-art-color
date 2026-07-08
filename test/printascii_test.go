package test

import (
	"bytes"
	"strings"
	"testing"

	"ascii-art/internal"
)

func TestEmptyLine(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		nil, nil,
		[]string{""},
		"../banners/standard.txt",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

func TestSingleWord(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		nil, nil,
		[]string{"Hi"},
		"../banners/standard.txt",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// TrimRight strips the trailing newline that Fprintln appends to each row,
	// otherwise Split produces a spurious empty element and len == 9.
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	if len(lines) != 8 {
		t.Errorf("expected 8 lines of ASCII art, got %d", len(lines))
	}
}

func TestNewlineSeparator(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		nil, nil,
		[]string{"Hi", "", "There"},
		"../banners/standard.txt",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// "Hi"     -> 8 lines
	// ""       -> 1 empty line
	// "There"  -> 8 lines
	//
	// Total: 17 lines
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	if len(lines) != 17 {
		t.Errorf("expected 17 lines (8 + 1 + 8), got %d", len(lines))
	}
}

func TestSpecificCharacter(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(
		&buf,
		nil, nil,
		[]string{"A"},
		"../banners/standard.txt",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	lines := strings.Split(output, "\n")

	// lines[4] is the 5th row (0-indexed); 'A' has its crossbar there.
	if !strings.Contains(lines[4], "_") {
		t.Errorf(
			"expected 5th line of 'A' to contain '_', got %q",
			lines[4],
		)
	}
}
