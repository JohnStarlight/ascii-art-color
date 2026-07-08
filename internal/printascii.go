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
	colors, parts []string,
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

		ranges := buildColorRanges(line, colors, parts)

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
				if color := colorAt(pos, ranges); color != "" {
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

type colorRange struct {
	color   string
	starts  []int
	partLen int
}

// buildColorRanges pairs each color with its substring occurrences in line,
// preserving flag order so earlier --color flags take priority on overlap.
func buildColorRanges(line string, colors, parts []string) []colorRange {
	var ranges []colorRange
	for i, color := range colors {
		if color == "" {
			continue

		}

		var part string
		if i < len(parts) {
			part = parts[i]
		}

		partLen := len(part)
		if part == "" {
			partLen = len(line)
		}

		ranges = append(ranges, colorRange{
			color:   color,
			starts:  findAll(line, part),
			partLen: partLen,
		})
	}
	return ranges
}

func colorAt(pos int, ranges []colorRange) string {
	for _, r := range ranges {
		if inColorRange(pos, r.starts, r.partLen) {
			return r.color
		}
	}
	return ""
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
