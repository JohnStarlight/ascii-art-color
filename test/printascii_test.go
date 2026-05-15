package test

// import (
// 	"bytes"
// 	"io"
// 	"os"
// 	"strings"
// 	"testing"

//	"ascii-art/internal"
//
// )
func forDeletion() {}

// // captureOutput "αρπάζει" ό,τι τυπώνεται στο stdout και το επιστρέφει ως string
// func captureOutput(f func()) string {
// 	r, w, _ := os.Pipe()
// 	old := os.Stdout
// 	os.Stdout = w

// 	f()

// 	w.Close()
// 	os.Stdout = old

// 	var buf bytes.Buffer
// 	io.Copy(&buf, r)
// 	return buf.String()
// }

// // Ελέγχει ότι κενή είσοδος δεν τυπώνει τίποτα
// func TestEmptyLine(t *testing.T) {
// 	output := captureOutput(func() {
// 		err := internal.PrintAscii([]string{""}, "../banners/standard.txt")
// 		if err != nil {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
// 	})

// 	if output != "" {
// 		t.Errorf("expected empty output, got %q", output)
// 	}
// }

// // Ελέγχει ότι η έξοδος για ένα string έχει ακριβώς 8 γραμμές ASCII art
// func TestSingleWord(t *testing.T) {
// 	output := captureOutput(func() {
// 		err := internal.PrintAscii([]string{"Hi"}, "../banners/standard.txt")
// 		if err != nil {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
// 	})

// 	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
// 	if len(lines) != 8 {
// 		t.Errorf("expected 8 lines of ASCII art, got %d", len(lines))
// 	}
// }

// // Ελέγχει ότι το \n χωρίζει σωστά σε δύο blocks ASCII art με μία κενή γραμμή
// func TestNewlineSeparator(t *testing.T) {
// 	output := captureOutput(func() {
// 		err := internal.PrintAscii([]string{"Hi", "", "There"}, "../banners/standard.txt")
// 		if err != nil {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
// 	})

// 	// "Hi" = 8 γραμμές, "" = 1 κενή γραμμή, "There" = 8 γραμμές → σύνολο 17
// 	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
// 	if len(lines) != 17 {
// 		t.Errorf("expected 17 lines (8 + 1 + 8), got %d", len(lines))
// 	}
// }

// // Ελέγχει ότι το A έχει το σωστό ASCII art (πέμπτη γραμμή του A στο standard έχει _)
// func TestSpecificCharacter(t *testing.T) {
// 	output := captureOutput(func() {
// 		err := internal.PrintAscii([]string{"A"}, "../banners/standard.txt")
// 		if err != nil {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
// 	})

// 	lines := strings.Split(output, "\n")
// 	if !strings.Contains(lines[4], "_") {
// 		t.Errorf("expected 5th line of 'A' to contain '_', got %q", lines[4])
// 	}
// }
