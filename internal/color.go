package internal

import (
	"fmt"
	"strconv"
	"strings"
)

func Palette(color string) (string, error) {
	R, G, B := 0, 0, 0
	color = strings.ToLower(color)
	if strings.HasPrefix(color, "rgb(") {
		part := strings.TrimPrefix(color, "rgb(")
		part = strings.TrimSuffix(part, ")")
		vals := strings.Split(part, ",")
		if len(vals) != 3 {
			return "", fmt.Errorf("invalid rgb format: expected rgb(R,G,B)")
		}
		var err error
		if R, err = strconv.Atoi(strings.TrimSpace(vals[0])); err != nil {
			return "", fmt.Errorf("invalid rgb red value")
		}
		if G, err = strconv.Atoi(strings.TrimSpace(vals[1])); err != nil {
			return "", fmt.Errorf("invalid rgb green value")
		}
		if B, err = strconv.Atoi(strings.TrimSpace(vals[2])); err != nil {
			return "", fmt.Errorf("invalid rgb blue value")
		}
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
