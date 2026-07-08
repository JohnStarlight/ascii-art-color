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

	err := internal.PrintAscii(&buf, nil, nil, []string{"Hi"}, "../banners/standard.txt")
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

	err = internal.PrintAscii(&buf, []string{color}, nil, []string{"Hi"}, "../banners/standard.txt")
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
	err = internal.PrintAscii(&buf, []string{color}, []string{"H"}, []string{"Hi"}, "../banners/standard.txt")
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
	err = internal.PrintAscii(&plainBuf, nil, nil, []string{"Hi"}, "../banners/standard.txt")
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
	if len(cfg.Colors) != 1 || cfg.Colors[0] != want {
		t.Errorf("Colors = %v, want [%q]", cfg.Colors, want)
	}
	if len(cfg.Parts) != 1 || cfg.Parts[0] != "" {
		t.Errorf("Parts = %v, want [\"\"]", cfg.Parts)
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

	if len(cfg.Parts) != 1 || cfg.Parts[0] != "ell" {
		t.Errorf("Parts = %v, want [%q]", cfg.Parts, "ell")
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
	if len(cfg.Parts) != 1 || cfg.Parts[0] != "" {
		t.Errorf("Parts = %v, want [\"\"]", cfg.Parts)
	}
	if cfg.Text != "Hello" {
		t.Errorf("Text = %q, want %q", cfg.Text, "Hello")
	}
}

func TestParseArgsTwoColors(t *testing.T) {
	cfg, err := internal.ParseArgs([]string{"prog", "--color=red", "--color=blue", "ell", "lo", "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantRed, _ := internal.Palette("red")
	wantBlue, _ := internal.Palette("blue")
	if len(cfg.Colors) != 2 || cfg.Colors[0] != wantRed || cfg.Colors[1] != wantBlue {
		t.Errorf("Colors = %v, want [%q %q]", cfg.Colors, wantRed, wantBlue)
	}
	if len(cfg.Parts) != 2 || cfg.Parts[0] != "ell" || cfg.Parts[1] != "lo" {
		t.Errorf("Parts = %v, want [ell lo]", cfg.Parts)
	}
	if cfg.Text != "Hello" {
		t.Errorf("Text = %q, want %q", cfg.Text, "Hello")
	}
}

func TestParseArgsTwoColorsWithBanner(t *testing.T) {
	cfg, err := internal.ParseArgs([]string{"prog", "--color=red", "--color=blue", "ell", "lo", "Hello", "shadow"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BannerPath != "banners/shadow.txt" {
		t.Errorf("BannerPath = %q, want %q", cfg.BannerPath, "banners/shadow.txt")
	}
	if cfg.Text != "Hello" {
		t.Errorf("Text = %q, want %q", cfg.Text, "Hello")
	}
}

func TestParseArgsTwoColorsMissingSubstring(t *testing.T) {
	_, err := internal.ParseArgs([]string{"prog", "--color=red", "--color=blue", "Hello"})
	if err == nil {
		t.Fatal("expected error when a substring is missing for one of two --color flags, got nil")
	}
}

func TestPrintAsciiTwoColors(t *testing.T) {
	var buf bytes.Buffer

	red, err := internal.Palette("red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	blue, err := internal.Palette("blue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// "Hello": "H" red, "lo" blue.
	err = internal.PrintAscii(&buf, []string{red, blue}, []string{"H", "lo"}, []string{"Hello"}, "../banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	for i, line := range lines {
		if !strings.Contains(line, red) {
			t.Errorf("line %d: expected line to contain red color code, got %q", i, line)
		}
		if !strings.Contains(line, blue) {
			t.Errorf("line %d: expected line to contain blue color code, got %q", i, line)
		}
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
