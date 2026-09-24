# CSV ⇄ Parquet Converter

[![Go Version](https://img.shields.io/badge/go-1.27+-blue?logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/dbunt1tled/parquet2csv.svg)](https://pkg.go.dev/github.com/dbunt1tled/parquet2csv)
[![Build Status](https://github.com/dbunt1tled/parquet2csv/workflows/Build/badge.svg)](https://github.com/dbunt1tled/parquet2csv/actions)
[![Release](https://img.shields.io/github/v/release/dbunt1tled/parquet2csv)](https://github.com/dbunt1tled/parquet2csv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/dbunt1tled/parquet2csv)](https://goreportcard.com/report/github.com/dbunt1tled/parquet2csv)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A CLI tool for **bidirectional conversion** between CSV and Apache Parquet, built in Go on the
Cobra CLI framework. It streams both directions a row at a time, so peak memory is set by the
tuning flags rather than by the size of the input file.

## Features

- 🔄 **Bidirectional conversion**: CSV → Parquet and Parquet → CSV
- 🌊 **Streaming**: one record in flight at a time in both directions; no whole-file buffering
- 📐 **Flat memory**: peak usage is bounded by `--row-group-size` and `--page-size`, not by row count
- 🗜️ **Compression**: UNCOMPRESSED, SNAPPY, GZIP, LZ4 and ZSTD
- ✂️ **Custom delimiters**: any single character, including multi-byte runes
- 📊 **Verbose statistics**: elapsed time and memory usage per run

> **No type inference.** Every Parquet column is written as `BYTE_ARRAY` / `UTF8`. Numbers and
> dates survive the round trip as text. See [Limitations](#limitations).

## Installation

### From source

```bash
git clone https://github.com/dbunt1tled/parquet2csv.git
cd parquet2csv
go build -o csv2parquet main.go
```

Requires **Go 1.27** or newer (`go.mod` targets `go 1.27.1`).

### Using go install

```bash
go install github.com/dbunt1tled/parquet2csv@latest
```

> ⚠️ The module is named `parquet2csv` while the command is `csv2parquet` — the repository was
> renamed and the import path was not. `go install` therefore places a binary called
> **`parquet2csv`** in `$GOBIN`. Rename it if you want the documented name:
> `mv "$(go env GOPATH)/bin/parquet2csv" "$(go env GOPATH)/bin/csv2parquet"`.

### Prebuilt binaries

Archives are published on the [releases page](https://github.com/dbunt1tled/parquet2csv/releases)
for linux, macOS and windows on amd64 and arm64 (no windows/arm64). The binary inside is named
`csv2parquet`.

## Usage

### Command tree

```
csv2parquet                       # root command; prints help, converts nothing
  ├── parquet <input.csv> [output]      # CSV → Parquet
  ├── csv <input.parquet> [output]      # Parquet → CSV
  ├── completion <shell>                # shell completion script
  └── help [command]                    # help for any command
```

**Subcommands are named after the output format, not the input.** `csv2parquet parquet` reads a
CSV, `csv2parquet csv` reads a Parquet file.

### How arguments are resolved

- The **input** argument is required, and its extension is checked strictly: `parquet` rejects
  anything that is not `.csv`, `csv` rejects anything that is not `.parquet`.
- The **output** argument is optional. Left out, the input's extension is swapped:
  `data.csv` → `data.parquet`, `data.parquet` → `data.csv`.
- When given, the output extension is **forced**, not validated. `out` becomes `out.parquet`,
  and `weird.txt` becomes `weird.txt.parquet` — only an already-correct extension is left alone.
- The output directory must exist and be writable, and **an existing output file is overwritten
  without a prompt.** Converting `data.parquet` with no output argument writes `data.csv`, which
  will silently replace the CSV it was made from.

### Flags — `csv2parquet parquet` (CSV → Parquet)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--compression` | `-c` | int | `0` | Codec, see the table below |
| `--delimiter` | `-d` | string | `,` | Input field delimiter; exactly one character |
| `--flush` | `-f` | int | `10000` | Rows between calls that turn buffered rows into pages |
| `--row-group-size` | `-r` | int | `8` | Target row group size, in MB |
| `--page-size` | `-p` | int | `1024` | Target data page size, in KB |
| `--verbose` | `-v` | bool | `false` | Print elapsed time and memory usage when done |
| `--help` | `-h` | bool | `false` | Show help for the subcommand |

### Flags — `csv2parquet csv` (Parquet → CSV)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--delimiter` | `-d` | string | `,` | Output field delimiter; exactly one character |
| `--flush` | `-f` | int | `10000` | Rows per read chunk, and the CSV writer's flush cadence |
| `--verbose` | `-v` | bool | `false` | Print elapsed time and memory usage when done |
| `--compression` | `-c` | int | `0` | **Accepted but ignored** — the flag is registered and never read |
| `--help` | `-h` | bool | `false` | Show help for the subcommand |

### Compression codecs

| Value | Codec |
|-------|-------|
| `0` | UNCOMPRESSED (default) |
| `1` | SNAPPY |
| `2` | GZIP |
| `5` | LZ4 |
| `6` | ZSTD |

Any other value is rejected up front: parquet-go silently compresses to nothing for a codec it
has no compressor for, which would write a file with a schema and no rows.

## All ways to run it

### CSV → Parquet

```bash
# Output name derived from the input: data.csv → data.parquet
csv2parquet parquet data.csv

# Explicit output path
csv2parquet parquet data.csv /tmp/out.parquet

# Output extension is forced, so this also writes out.parquet
csv2parquet parquet data.csv out

# Semicolon-delimited input
csv2parquet parquet data.csv -d ';'

# Tab-delimited input (quote the escape so the shell expands it).
# The file still has to be named .csv — the extension check is strict.
csv2parquet parquet tabbed.csv -d $'\t'

# Multi-byte delimiter
csv2parquet parquet data.csv -d '→'

# Snappy compression, with statistics
csv2parquet parquet data.csv -c 1 -v

# Zstd, the smallest output of the supported codecs
csv2parquet parquet data.csv -c 6

# Larger row groups, better for analytical readers, at the cost of peak memory
csv2parquet parquet data.csv -r 128

# Everything at once, long form
csv2parquet parquet big_file.csv big_file.parquet \
  --delimiter ';' \
  --compression 2 \
  --flush 50000 \
  --row-group-size 64 \
  --page-size 2048 \
  --verbose
```

### Parquet → CSV

```bash
# Output name derived from the input: data.parquet → data.csv
csv2parquet csv data.parquet

# Explicit output path
csv2parquet csv data.parquet /tmp/out.csv

# Pipe-delimited output
csv2parquet csv analytics.parquet analytics.csv -d '|'

# Smaller read chunks and flush interval
csv2parquet csv analytics.parquet -f 1000 -v
```

A Parquet file with zero rows still carries a schema, so it produces a CSV containing the header
row and nothing else — not an empty file.

### Help and completion

```bash
csv2parquet                       # root help; no conversion happens
csv2parquet --help                # same
csv2parquet parquet --help        # CSV → Parquet help
csv2parquet csv --help            # Parquet → CSV help
csv2parquet help parquet          # equivalent to the above

csv2parquet completion bash       # completion script for bash
csv2parquet completion zsh        # …zsh
csv2parquet completion fish       # …fish
csv2parquet completion powershell # …powershell
```

There is **no `--version` flag**. The release build injects `main.version`, `main.commit` and
`main.date` via ldflags, but `main.go` declares no such variables, so the values go nowhere.

### Exit behaviour

A failure prints `Error: <message>` and exits with status 1. The common ones:

| Message | Cause |
|---------|-------|
| `file is not csv file` / `file is not parquet file` | Input extension does not match the subcommand |
| `input file X not exist: ...` | Input path does not resolve |
| `input file X is empty: no header row` | CSV had no records at all |
| `delimiter must be a single character, got "..."` | `--delimiter` was empty or longer than one rune |
| `flush must be at least 1, got N` | Non-positive `--flush` |
| `row-group-size must be at least 1 MB, got N` | Non-positive `--row-group-size` |
| `page-size must be at least 1 KB, got N` | Non-positive `--page-size` |
| `unsupported compression N, want one of [0 1 2 5 6]` | Unknown codec |
| `Path doesn't exist. X` / `Path isn't a directory. X` | Output directory missing or not a directory |

## Tuning memory and speed

Peak memory for CSV → Parquet is set by two flags rather than by the size of the input:

| Flag | What it bounds |
|------|----------------|
| `--row-group-size` | Pages held in memory before a row group is written out |
| `--page-size` | The page count, and with it the column and offset indexes kept until the footer is written |

`--row-group-size` is the dominant dial: pages accumulate in memory until a row group is written
out, so raising it raises peak memory directly. `--page-size` is second order and pulls in two
opposing directions — parquet-go keeps one column index and one offset index entry **per page**
until the file is finalised, so a small page size makes peak memory creep up with the length of
the file, while the writer buffers rows in proportion to page size times column count, so a large
one costs a fixed buffer instead. The default of 1024 KB matches Arrow/parquet-cpp and is the
better trade from a few million rows up; on short inputs parquet-go's 8 KB default is slightly
leaner.

Measured on a synthetic CSV of 25 string columns, uncompressed output:

| Input | Peak RSS | Time |
|-------|----------|------|
| 644 MB / 2M rows | 125 MB | 4.3 s |
| 1.3 GB / 4M rows | 129 MB | 8.8 s |

Doubling the input leaves peak memory essentially unchanged. `--row-group-size` is the dial that
moves it — on the 4M-row file above, with SNAPPY:

| `--row-group-size` | Peak RSS | Output size |
|--------------------|----------|-------------|
| `8` (default) | 69 MB | 224 MB |
| `32` | 118 MB | 224 MB |
| `128` | 314 MB | 224 MB |

The output is the same size either way here, so lowering it costs nothing on disk; what larger
row groups buy is read-side efficiency in analytical engines, which prefer fewer, bigger groups.

Cross-checked on a real 1.08 GiB / 4M-row / 37-column CSV, uncompressed: 87 MB at the defaults,
132 MB at `-r 32`, 348 MB at `-r 128`, and 110 MB at `-r 8 -p 8`. The default of 8 MB was chosen
from those numbers: it costs nothing in conversion time against 32 MB, and going below it buys
nothing — with `GOGC=20` to remove the garbage collector from the comparison, `-r 8`, `-r 4` and
`-r 1` all land within 1 MB of each other.

`--flush` does **not** control row group boundaries. It only sets how often buffered rows are
converted to pages, and rarely needs changing.

## Limitations

- **No type inference.** Every column is written as `BYTE_ARRAY` / `UTF8`, whatever the CSV
  contained. A round trip preserves the text exactly, but numeric and date columns arrive in
  Parquet as strings.
- **Column names must be usable in a Parquet tag.** A name containing a quote or a comma is
  rejected, as are empty and duplicate names.
- **LZ4 (`-c 5`) writes the deprecated `LZ4` codec**, not `LZ4_RAW`. Some readers, DuckDB among
  them, refuse such files. Prefer `-c 6` (ZSTD) or `-c 1` (SNAPPY) for interchange.
- **`--compression` on the `csv` subcommand does nothing.** It is registered and never read.
- **Existing output files are overwritten without confirmation.**

## Dependencies

- **CLI framework**: `github.com/spf13/cobra v1.10.2`
- **Parquet**: `github.com/xitongsys/parquet-go v1.6.2` and `parquet-go-source`
- **Error wrapping**: `github.com/pkg/errors v0.9.1`
- **String utilities**: `github.com/iancoleman/strcase v0.3.0`
- **Dynamic structs**: `github.com/ompluscator/dynamic-struct v1.4.0`

## Development

### Project structure

```
├── cmd/                   # Cobra CLI commands
│   ├── root.go            # Root command definition
│   ├── csv2parquet.go     # CSV → Parquet conversion
│   └── parquet2csv.go     # Parquet → CSV conversion
├── internal/
│   ├── file/              # Streaming CSV reader, CSV writer, path checks
│   ├── helper/            # Conversion and formatting utilities
│   └── schema/            # Dynamic struct and Parquet tag generation
└── main.go                # Application entry point
```

### Building

```bash
go build -o csv2parquet main.go                              # plain build
GOGC=150 CGO_ENABLED=0 go build -ldflags "-s -w" \
  -o csv2parquet main.go                                     # release flags, as used by CI
```

The `Makefile` targets are currently broken (`@GOGC=150 @go build` expands to
`@go: command not found`); use `go build` directly.

### Testing

```bash
go test ./...                 # all tests
go test -v ./...              # verbose
go test -race ./...           # with the race detector
go test -bench . ./...        # benchmarks
```

### Linting

```bash
golangci-lint run             # config: .golangci.yml
```

`golangci-lint` will refuse to start if its binary was built with an older Go than the module
targets (`the Go language version (go1.25) ... is lower than the targeted Go version (1.27.1)`).
Rebuild it against Go 1.27; `go vet ./...` is a partial fallback.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Parquet handling by [xitongsys/parquet-go](https://github.com/xitongsys/parquet-go)
- CLI powered by [spf13/cobra](https://github.com/spf13/cobra)
- Dynamic schema structs via [ompluscator/dynamic-struct](https://github.com/ompluscator/dynamic-struct)
