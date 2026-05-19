package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ascii-art/internal"
)

func main() {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Println("Error: invalid usage")
		fmt.Println("Usage: go run ./cmd --flag:Color [optional]\"your-colorizer-here\" \"your-text-here\"")
		os.Exit(1)
	}

	flag := os.Args[1]
	part := ""
	if len(os.Args) != 3 {
		part = os.Args[len(os.Args)-2]
	}
	text := os.Args[len(os.Args)-1]

	// Only printable ASCII characters (32–126) are supported.
	// Reject anything outside that range (e.g. accented letters, emoji).
	for _, r := range text {
		if r < 32 || r > 126 {
			fmt.Printf("Error: invalid character %q (only printable ASCII is supported)\n", r)
			os.Exit(1)
		}
	}

	fmt.Println("In which style would you like that?")
	fmt.Println("1 = Standard")
	fmt.Println("2 = Shadow")
	fmt.Println("3 = Thinkertoy")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error: failed to read input: %v\n", err)
		os.Exit(1)
	}
	input = strings.TrimSpace(input)

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > 3 {
		fmt.Println("Error: invalid choice (must be 1, 2 or 3)")
		os.Exit(1)
	}

	banners := []string{"banners/standard.txt", "banners/shadow.txt", "banners/thinkertoy.txt"}
	filename := banners[choice-1]

	flagParts := strings.Split(flag, "=")
	if len(flagParts) != 2 || flagParts[0] != "--color" {
		// error msg not proper flag
		fmt.Println("Usage: go run . [OPTION] [STRING]  EX: go run . --color=<color> <substring to be colored> \"something\"")
		os.Exit(1)
	}

	// Split on the literal two-character sequence "\n" (backslash + n),
	// which is how multi-line input is passed from the command line.
	// e.g. "Hello\nThere" becomes ["Hello", "There"].
	lines := strings.Split(text, "\\n")

	color, err := pallete(flagParts[1])

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = internal.PrintAscii(color, part, lines, filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func pallete(color string) (string, error) {
	R, G, B := 0, 0, 0
	color = strings.ToLower(color)
	if strings.HasPrefix(color, "rgb(") {
		// extract the rgb values
		part := strings.TrimPrefix(color, "rgb(")
		part = strings.TrimSuffix(part, ")")
		ColorValues := strings.Split(part, ",")
		// split the color values
		R, _ = strconv.Atoi(ColorValues[0])
		G, _ = strconv.Atoi(ColorValues[1])
		B, _ = strconv.Atoi(ColorValues[2])
	} else {
		switch color {
		case "red":
			R = 255
		case "green":
			G = 255
		case "blue":
			B = 255
		case "yellow":
			R, G = 255, 255
		case "magenta", "purple":
			R, B = 255, 255
		case "cyan":
			G, B = 255, 255
		default:
			return "", fmt.Errorf("unknown color: %s", color)
		}
	}

	return fmt.Sprintf("\033[38;2;%d;%d;%dm", R, G, B), nil
}
