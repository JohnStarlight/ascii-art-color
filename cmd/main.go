package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ascii-art/internal"
)

// run exists so that main stays minimal and all errors surface in one place
// instead of calling os.Exit at multiple points throughout the code.
func run() error {
	config, err := internal.ParseArgs(os.Args)
	if err != nil {
		return err
	}

	var writer io.Writer = os.Stdout

	if config.OutputFile != "" {
		file, err := internal.PrepareOutputFile(config.OutputFile)
		if err != nil {
			return err
		}

		defer file.Close()

		writer = file
	}

	// Split on the literal two-character sequence \n as typed on the CLI.
	lines := strings.Split(config.Text, "\\n")

	if err := internal.PrintAscii(
		writer,
		config.Color,
		config.Part,
		lines,
		config.BannerPath,
	); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		os.Exit(1)
	}
}
