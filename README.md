# Dircard

[English](README.md) | [日本語](README.ja.md)

Dircard is a lightweight CLI that shows directory notes from a `.dircard` file.

It is designed for people who move across many projects and want local context to appear automatically in the terminal.

## Features

- Show notes from the nearest `.dircard` in the current or parent directories.
- Limit output size, line range, and search depth.
- JSON output support for scripting.
- Shell integration for `bash`, `zsh`, and `pwsh`.
- Safe install and uninstall flow for shell hooks.

## Installation

Install with Go:

```bash
go install github.com/dircard/dircard@latest
```

## Quick Start

1. Create a `.dircard` file in your project directory.
2. Add any notes you want to see when entering that directory.
3. Check the output manually:

```bash
dircard show
```

4. Enable shell integration:

```bash
dircard install bash
# or
dircard install zsh
# or
dircard install pwsh
```

Then restart your shell (or source your rc file) to apply changes.

5. Check automatic display on directory change:

```bash
cd your/project
# The contents of .dircard will be shown on directory change
```

## Commands

### Show notes

```bash
dircard show
dircard show --full
dircard show --path
dircard show --json
dircard show --depth 3 --start 10 --lines 20
```

- `dircard show` searches for the nearest `.dircard` from the current directory upward.
- `--full` shows the full contents of the file.
- `--path` shows the path to the `.dircard` file.
- `--json` prints structured output for scripts.
- `--depth`, `--start`, and `--lines` control search depth and output range.

### Install shell hook

```bash
dircard install [bash|zsh|pwsh]
dircard install pwsh --force
```

`--force` updates only the dircard hook block and does not overwrite unrelated file contents.

### Uninstall shell hook

```bash
dircard uninstall [bash|zsh|pwsh]
dircard uninstall --force
```

Without a shell argument, `uninstall` removes hooks from all supported shells.

## Development

Build locally:

```bash
go build ./...
```

Run the CLI during development:

```bash
go run . show
```

Inspect available commands and flags:

```bash
go run . --help
go run . show --help
```

## Author

yhotta240 [https://github.com/yhotta240](https://github.com/yhotta240)

## License

MIT
