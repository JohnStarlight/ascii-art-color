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
	if len(os.Args) != 4 {
		fmt.Println("Error: invalid usage")
		fmt.Println("Usage: go run ./cmd --flag:Color \"your-colorizer-here\" \"your-text-here\"") // will be immersive
		os.Exit(1)
	}

	flag := os.Args[1]
	part := os.Args[2]
	text := os.Args[3]

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

	color := strings.Split(flag, "=")

	// Split on the literal two-character sequence "\n" (backslash + n),
	// which is how multi-line input is passed from the command line.
	// e.g. "Hello\nThere" becomes ["Hello", "There"].
	lines := strings.Split(text, "\\n")

	err = internal.PrintAscii(color[1], part, lines, filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
