package test

import (
	"bytes"
	"strings"
	"testing"

	"ascii-art/internal"
)

func expectedRenderedLines(lines []string) int {
	// Strip trailing empty segments to match TrimRight behaviour in assertions.
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	total := 0

	for i, line := range lines {
		if line == "" {
			if i > 0 {
				total++
			}
			continue
		}

		total += 8 // each ASCII-art character occupies 8 visual rows
	}

	return total
}

func renderStandard(input string) string {
	var buf bytes.Buffer

	lines := strings.Split(input, "\\n") // literal \n as typed on the CLI, not an actual newline

	err := internal.PrintAscii(&buf, "", "", lines, "../banners/standard.txt")
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func TestAuditExamplesNonEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		mustContain []string
	}{
		{
			name:        "hello",
			input:       "hello",
			mustContain: []string{"| |__     ___"},
		},
		{
			name:        "HELLO",
			input:       "HELLO",
			mustContain: []string{"|  ____|", "| |__"},
		},
		{
			name:        "HeLlo HuMaN",
			input:       "HeLlo HuMaN",
			mustContain: []string{"| \\  / |", "| \\ | |"},
		},
		{
			name:        "1Hello 2There",
			input:       "1Hello 2There",
			mustContain: []string{"|___ \\", "|__   __|"},
		},
		{
			name:        "Hello\\nThere",
			input:       "Hello\\nThere",
			mustContain: []string{"|__   __|", "|_|  |_|"},
		},
		{
			name:        "Hello\\n\\nThere",
			input:       "Hello\\n\\nThere",
			mustContain: []string{"|__   __|", "\n\n _______"},
		},
		{
			name:        "{Hello & There #}",
			input:       "{Hello & There #}",
			mustContain: []string{"\\___/\\/", "_| || |_"},
		},
		{
			name:        "hello There 1 to 2!",
			input:       "hello There 1 to 2!",
			mustContain: []string{"|___ \\", "|  _ \\"},
		},
		{
			name:        "MaD3IrA&LiSboN",
			input:       "MaD3IrA&LiSboN",
			mustContain: []string{"|_____/  |____/", "|_.__/"},
		},
		{
			name:        "1a\\\"#FdwHywR&/()=",
			input:       "1a\\\"#FdwHywR&/()=",
			mustContain: []string{"|______|", "__/ /"},
		},
		{
			name:        "{|}~",
			input:       "{|}~",
			mustContain: []string{"/\\/|", "|_|"},
		},
		{
			name:        "[\\]^_ 'a",
			input:       "[\\]^_ 'a",
			mustContain: []string{"|______|", "\\__,_|"},
		},
		{
			name:        "RGB",
			input:       "RGB",
			mustContain: []string{"_____   ____", "| |__) |"},
		},
		{
			name:        ":;<=>?@",
			input:       ":;<=>?@",
			mustContain: []string{"|______|", "\\____/"},
		},
		{
			name:        "\\!\" #$%&'()*+,-./",
			input:       "\\!\" #$%&'()*+,-./",
			mustContain: []string{"|______|", "\\_\\ /_/"},
		},
		{
			name:        "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			input:       "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			mustContain: []string{"|  _ \\   / ____|", "|_____/     |_|"},
		},
		{
			name:        "abcdefghijklmnopqrstuvwxyz",
			input:       "abcdefghijklmnopqrstuvwxyz",
			mustContain: []string{"__ _  | |__", "|_.__/   \\___|"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := renderStandard(tc.input)

			if output == "" {
				t.Fatalf("expected non-empty output for %q", tc.input)
			}

			lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
			wantLines := expectedRenderedLines(strings.Split(tc.input, "\\n"))

			if len(lines) != wantLines {
				t.Fatalf(
					"expected %d rendered lines, got %d",
					wantLines,
					len(lines),
				)
			}

			for _, fragment := range tc.mustContain {
				if !strings.Contains(output, fragment) {
					t.Fatalf(
						"expected output to contain %q",
						fragment,
					)
				}
			}
		})
	}
}
