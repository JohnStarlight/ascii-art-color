# ASCII Art Color CLI

A Go command-line tool that renders text as ASCII art with optional terminal color highlighting.

## What This Project Does

- Takes a required `--color=<color>` flag and a text argument
- Optionally takes a substring argument to color only that part of the output
- Converts each character to ASCII art (8 lines tall)
- Lets you choose one of 3 styles interactively

## Quick Start

From the project root:

```bash
go run ./cmd --color=red "Hello"
```

You will then be prompted to choose a style:

```text
In which style would you like that?
1 = Standard
2 = Shadow
3 = Thinkertoy
```

Enter `1`, `2`, or `3` and press Enter.

## Usage

Run from the project root (not inside `cmd/`):

```text
go run ./cmd --color=<color> [substring] "your text"
```

| Argument      | Required | Description                                              |
|---------------|----------|----------------------------------------------------------|
| `--color=<c>` | Yes      | Color to apply (see supported colors below)              |
| `[substring]` | No       | If provided, only this part of the text is colored       |
| `"your text"` | Yes      | The text to render as ASCII art                          |

### Color only a substring

```bash
go run ./cmd --color=green "Hello" "Hello There"
```

This colors only the word `Hello` in green, leaving `There` uncolored.

### Color the entire text

```bash
go run ./cmd --color=cyan "Hello There"
```

Omitting the substring argument colors the entire text.

## Supported Colors

Named colors:

- `red`
- `green`
- `blue`
- `yellow`
- `cyan`
- `magenta` / `purple`

RGB values:

```bash
go run ./cmd --color=rgb(255,128,0) "Hello"
```

## Multi-line Output

Use the literal `\n` sequence inside the string to produce multiple lines:

```bash
go run ./cmd --color=blue "Hello\nThere"
```

## Input Rules

- Exactly 2 or 3 arguments are required (flag + text, or flag + substring + text)
- Only printable ASCII characters are accepted (`32` to `126`)
- Non-ASCII characters (e.g. accented letters, emoji) are rejected

## Run Tests

```bash
go test ./...
```

With verbose output:

```bash
go test ./... -v
```

Force a fresh run (no cache):

```bash
go test ./... -v -count=1
```

## Project Structure

- [cmd/main.go](cmd/main.go) — CLI entrypoint, argument parsing, color flag handling
- [internal/printascii.go](internal/printascii.go) — ASCII rendering and colorization logic
- [banners/](banners/) — banner template files (`standard.txt`, `shadow.txt`, `thinkertoy.txt`)
- [test/printascii_test.go](test/printascii_test.go) — core unit tests
- [test/audit_examples_test.go](test/audit_examples_test.go) — audit/instruction sample tests

## License

MIT. See `LICENSE`.
