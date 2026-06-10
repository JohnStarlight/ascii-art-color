package internal

import (
	"fmt"
	"os"
)

func PrepareOutputFile(filename string) (*os.File, error) {
	// os.Create truncates any existing file at this path.
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create output file: %w",
			err,
		)
	}

	return file, nil
}
