# PRD - ascii-art

Keep this short. The project subject + audit cases are the source of truth.

---

## 1. Problem Statement

This CLI receives text as input and prints the same text as ASCII-art using banner templates.

---

## 2. CLI Contract

Supported forms:

```
go run ./cmd "<text>"
go run ./cmd "<text>" <banner>
go run ./cmd --output=<filename> "<text>" <banner>
go run ./cmd --color=<color> "<text>"
go run ./cmd --color=<color> "<substring>" "<text>"
go run ./cmd --color=<color> "<substring>" "<text>" <banner>
```

- `<text>`: required, must be in quotes for multi-word input
- `<banner>`: optional positional arg; defaults to `standard`; accepted values: `standard`, `shadow`, `thinkertoy`
- `--output=<filename>`: optional flag; when present, writes ASCII art to file instead of stdout; banner becomes required
- `--color=<color>`: optional flag; colors the rendered output (see 3.5)
- `<substring>`: optional positional arg, only meaningful with `--color`; the portion of `<text>` to color. If omitted, the whole text is colored.
- `--color`/`--output` flags can appear in any order/position among the arguments
- Invalid `--output` format (missing `=` or empty filename): print error and exit with non-zero status
- Invalid `--color` format (missing `=`, empty value, or unknown color): print usage message and exit with non-zero status
- Invalid chars (outside ASCII 32-126):
  - print clear error message
  - exit with non-zero status
- Banner file read failure:
  - print `Could not open file: ...`
  - exit with non-zero status

---

## 3. Rendering Functions

Documented in plain language with examples.

### 3.1 Style Selection

- Banner is passed as a positional CLI argument (default: `standard`).
- Accepted values: `standard`, `shadow`, `thinkertoy` (case-insensitive).
- Unsupported banner name: print error and exit with non-zero status.

Example:
- `go run ./cmd "Hello" shadow` -> prints `Hello` in `shadow` font.

### 3.4 File Output

- When `--output=<filename>` is provided, ASCII art is written to that file instead of stdout.
- File is created if it does not exist; overwritten if it does.
- On file creation failure: print error and exit with non-zero status.
- Banner argument is required when `--output` is used.

Example:
- `go run ./cmd --output=result.txt "Hello" standard` -> writes `Hello` in `standard` font to `result.txt`.

### 3.2 Newline Handling

- Input is split on literal `\n`.
- Each segment is rendered as a separate 8-line block.
- Empty segments produce an empty line separator (except leading empty segment in current behavior).

Examples:
- `Hello\nThere` -> two rendered blocks with one blank line between.
- `Hello\n\nThere` -> two rendered blocks with two separators (one empty rendered segment in between).

### 3.3 Character-to-Glyph Mapping

- Each printable ASCII character maps to one 8-line glyph from the banner file.
- For each output row (1..8), the program concatenates the correct row from each character glyph.

Example:
- `A` -> prints 8 lines that correspond to ASCII code 65 in the selected banner.

### 3.5 Color Highlighting

- `--color=<color>` colors the rendered ASCII art using ANSI 24-bit (`\033[38;2;R;G;Bm`) escape codes.
- Accepted color names: `red`, `green`, `blue`, `yellow`, `magenta`/`purple`, `cyan`.
- Custom colors via `rgb(R,G,B)`, e.g. `"--color=rgb(255,128,0)"` (quotes required, see below).
- An optional `<substring>` positional argument selects which part of `<text>` is colored; all occurrences of the substring are colored.
- If `<substring>` is omitted, the entire rendered text is colored.
- Unknown color name, malformed `rgb(...)`, or malformed `--color` flag (missing `=` or empty value): print usage message and exit with non-zero status.
- Shell note: `rgb(...)` contains parentheses, which are special characters in most shells. The whole flag must be quoted, e.g. `"--color=rgb(0,200,255)"`.

Examples:
- `go run ./cmd --color=red "Hello"` -> entire `Hello` rendered in red.
- `go run ./cmd --color=red "ell" "Hello"` -> only the `ell` substring rendered in red.
- `go run ./cmd "--color=rgb(0,200,255)" "Hello" shadow` -> entire `Hello` rendered in custom RGB color, `shadow` banner.

---

## 4. Non-Goals

- Unicode support beyond printable ASCII
- GUI/Web interface
- Custom user-uploaded fonts
- Rich text features beyond single-color highlighting (alignment, animation, gradients, multiple colors per call)

---

## 5. Acceptance Criteria

### Audit Cases

- [ ] `hello` prints expected ASCII-art output in selected style.
- [ ] Mixed case and spaces (for example `HeLlo HuMaN`) render correctly.
- [ ] Special characters (for example `{|}~` and punctuation sets) render correctly.
- [ ] Literal `\n` and `\n\n` create correct multi-line block separation.

### Extra Golden Tests

- [ ] Empty input behavior is defined and tested (`""`).
- [ ] Invalid character input (for example emoji) returns clear error and non-zero exit.
- [ ] Invalid style choice (non-numeric or out of range) returns clear error and non-zero exit.

---

## 6. Architecture

- We choose: Pipeline
- Because: input validation, banner selection, splitting, and rendering happen in a simple sequential flow.
- Tradeoffs we accept: less flexibility than a full parser/FSM, but easier to read and maintain.

Sketch (high-level):

`CLI args -> parse & validate -> load banner -> split by \n -> render each line -> stdout or file`

---

## 7. Milestones

1. Stabilize expected behavior for empty/newline edge cases.
2. Ensure functional output matches subject examples for key inputs.
3. Add/clean tests for validation, multiline behavior, and special chars.
4. Final docs pass (README usage + examples aligned with real behavior).

---

## 8. Risks / Open Questions

- Should final audit target exact `go run .` contract or keep `go run ./cmd`?

---

## 9. Known Limitations

- Rendering is limited to printable ASCII (`32-126`).
- Input with Unicode symbols (for example Greek characters or emoji) is not supported.
- When `--output` is used, banner argument is required (no fallback to default).
- Only one color can be applied per invocation; the colored substring is matched literally (no regex/wildcards).
- ANSI color codes are written even when output is redirected to a file via `--output` (no automatic stripping for non-terminal targets).
