package test

import (
	"bytes"
	"strings"
	"testing"

	"ascii-art/internal"
)

func TestPaletteNamedColors(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"red", "\033[38;2;255;0;0m"},
		{"green", "\033[38;2;0;255;0m"},
		{"blue", "\033[38;2;0;0;255m"},
		{"yellow", "\033[38;2;255;255;0m"},
		{"magenta", "\033[38;2;255;0;255m"},
		{"purple", "\033[38;2;255;0;255m"},
		{"cyan", "\033[38;2;0;255;255m"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := internal.Palette(tc.name)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Palette(%q) = %q, want %q", tc.name, got, tc.want)
			}
		})
	}
}

func TestPaletteRGB(t *testing.T) {
	got, err := internal.Palette("rgb(10,20,30)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "\033[38;2;10;20;30m"
	if got != want {
		t.Errorf("Palette(rgb(10,20,30)) = %q, want %q", got, want)
	}
}

func TestPaletteUnknownColor(t *testing.T) {
	_, err := internal.Palette("notacolor")
	if err == nil {
		t.Fatal("expected error for unknown color, got nil")
	}
}

func TestPaletteInvalidRGB(t *testing.T) {
	_, err := internal.Palette("rgb(1,2)")
	if err == nil {
		t.Fatal("expected error for malformed rgb(), got nil")
	}
}

func TestPrintAsciiNoColor(t *testing.T) {
	var buf bytes.Buffer

	err := internal.PrintAscii(&buf, "", "", []string{"Hi"}, "../banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(buf.String(), "\033[") {
		t.Errorf("expected no ANSI escape codes, got %q", buf.String())
	}
}

func TestPrintAsciiWholeStringColored(t *testing.T) {
	var buf bytes.Buffer

	color, err := internal.Palette("red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = internal.PrintAscii(&buf, color, "", []string{"Hi"}, "../banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	for i, line := range lines {
		if !strings.Contains(line, color) || !strings.Contains(line, "\033[0m") {
			t.Errorf("line %d: expected line to be wrapped in color codes, got %q", i, line)
		}
	}
}

func TestPrintAsciiPartialColor(t *testing.T) {
	var buf bytes.Buffer

	color, err := internal.Palette("red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// "Hi" with only "H" colored.
	err = internal.PrintAscii(&buf, color, "H", []string{"Hi"}, "../banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	for i, line := range lines {
		if !strings.Contains(line, color) {
			t.Errorf("line %d: expected line to contain color code for 'H', got %q", i, line)
		}
		if !strings.Contains(line, "\033[0m") {
			t.Errorf("line %d: expected line to contain reset code, got %q", i, line)
		}
	}

	// Without coloring, the same input should produce different (uncolored) output.
	var plainBuf bytes.Buffer
	err = internal.PrintAscii(&plainBuf, "", "", []string{"Hi"}, "../banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output == plainBuf.String() {
		t.Errorf("expected colored output to differ from plain output")
	}
}

func TestParseArgsColorWholeString(t *testing.T) {
	cfg, err := internal.ParseArgs([]string{"prog", "--color=red", "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want, _ := internal.Palette("red")
	if cfg.Color != want {
		t.Errorf("Color = %q, want %q", cfg.Color, want)
	}
	if cfg.Part != "" {
		t.Errorf("Part = %q, want empty", cfg.Part)
	}
	if cfg.Text != "Hello" {
		t.Errorf("Text = %q, want %q", cfg.Text, "Hello")
	}
}

func TestParseArgsColorWithSubstring(t *testing.T) {
	cfg, err := internal.ParseArgs([]string{"prog", "--color=red", "ell", "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Part != "ell" {
		t.Errorf("Part = %q, want %q", cfg.Part, "ell")
	}
	if cfg.Text != "Hello" {
		t.Errorf("Text = %q, want %q", cfg.Text, "Hello")
	}
}

func TestParseArgsColorWithBanner(t *testing.T) {
	cfg, err := internal.ParseArgs([]string{"prog", "--color=red", "Hello", "shadow"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BannerPath != "banners/shadow.txt" {
		t.Errorf("BannerPath = %q, want %q", cfg.BannerPath, "banners/shadow.txt")
	}
	if cfg.Part != "" {
		t.Errorf("Part = %q, want empty", cfg.Part)
	}
	if cfg.Text != "Hello" {
		t.Errorf("Text = %q, want %q", cfg.Text, "Hello")
	}
}

func TestParseArgsColorMissingEquals(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "--color", "Hello"})
	if err == nil {
		t.Fatal("expected error for --color without '=', got nil")
	}

	want := "Usage: go run . [OPTION] [STRING]\n\n" +
		"EX: go run . --color=<color> <substring to be colored> \"something\""
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestParseArgsColorEmptyValue(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "--color=", "Hello"})
	if err == nil {
		t.Fatal("expected error for --color= with empty value, got nil")
	}
}

func TestParseArgsUnknownColor(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "--color=notacolor", "Hello"})
	if err == nil {
		t.Fatal("expected error for unknown color, got nil")
	}
}
