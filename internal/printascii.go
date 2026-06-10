package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	asciiStart = 32 // ASCII value of ' ', the first printable character.

	charHeight = 8

	// 8 visual rows + 1 blank separator line per character in the banner file.
	linesPerChar = 9

	// 95 printable ASCII characters (32–126) × 9 lines each.
	expectedNewlines = 855
)

func PrintAscii(
	writer io.Writer,
	color, part string,
	lines []string,
	filename string,
) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(
			"could not open banner file: %w",
			err,
		)
	}

	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")

	if strings.Count(normalized, "\n") != expectedNewlines {
		return fmt.Errorf("banner file is corrupt or invalid")
	}

	bannerLines := strings.Split(normalized, "\n")

	for i, line := range lines {
		if line == "" {
			if i > 0 {
				fmt.Fprintln(writer)
			}
			continue
		}

		starts := findAll(line, part)
		partLen := len(part)
		if color != "" && part == "" {
			partLen = len(line)
		}

		for row := 1; row <= charHeight; row++ {
			var sb strings.Builder

			for pos, r := range line {
				index := (int(r)-asciiStart)*linesPerChar + row

				if index >= len(bannerLines) {
					return fmt.Errorf(
						"character %q is out of supported range in banner",
						r,
					)
				}

				segment := bannerLines[index]
				if color != "" && inColorRange(pos, starts, partLen) {
					sb.WriteString(color + segment + "\033[0m")
				} else {
					sb.WriteString(segment)
				}
			}

			fmt.Fprintln(writer, sb.String())
		}
	}

	return nil
}

func findAll(line, part string) []int {
	if part == "" {
		return []int{0}
	}
	var starts []int
	offset := 0
	for {
		idx := strings.Index(line[offset:], part)
		if idx == -1 {
			break
		}
		starts = append(starts, offset+idx)
		offset += idx + len(part)
	}
	return starts
}

func inColorRange(pos int, starts []int, partLen int) bool {
	for _, start := range starts {
		if pos >= start && pos < start+partLen {
			return true
		}
	}
	return false
}
