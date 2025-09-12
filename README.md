# CSV ⇄ Parquet Converter

[![Go Version](https://img.shields.io/badge/go-1.25+-blue?logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/dbunt1tled/parquet2csv.svg)](https://pkg.go.dev/github.com/dbunt1tled/parquet2csv)
[![Build Status](https://github.com/dbunt1tled/parquet2csv/workflows/Build/badge.svg)](https://github.com/dbunt1tled/parquet2csv/actions)
[![Release](https://img.shields.io/github/v/release/dbunt1tled/parquet2csv)](https://github.com/dbunt1tled/parquet2csv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/dbunt1tled/parquet2csv)](https://goreportcard.com/report/github.com/dbunt1tled/parquet2csv)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A fast, reliable CLI tool for **bidirectional conversion** between CSV and Apache Parquet formats. Built in Go with Cobra CLI framework, it's designed for data workflows that need efficient, schema-aware columnar storage with support for both directions of conversion.

## Features

- 🔄 **Bidirectional conversion**: CSV ↔ Parquet
- ⚡ **High performance**: Batch processing with configurable flush intervals
- 🗜️ **Compression support**: Multiple compression algorithms
- 🎯 **Schema-aware**: Automatic schema detection and type inference
- 📊 **Verbose statistics**: Runtime performance and memory usage reporting
- 🛠️ **Flexible CLI**: Powered by Cobra with intuitive subcommands

## Dependencies

- **Cobra CLI Framework**: `github.com/spf13/cobra v1.10.1`
- **Parquet Processing**: `github.com/xitongsys/parquet-go v1.6.2`
- **High-Performance JSON**: `github.com/bytedance/sonic v1.14.1`
- **Error Handling**: `github.com/pkg/errors v0.9.1`
- **String Utilities**: `github.com/iancoleman/strcase v0.3.0`
- **Dynamic Structs**: `github.com/ompluscator/dynamic-struct v1.4.0`

## Installation

### From Source
```bash
git clone https://github.com/dbunt1tled/parquet2csv.git
cd parquet2csv
go build -o csv2parquet main.go
```

### Using Go Install
```bash
go install github.com/dbunt1tled/parquet2csv@latest
```

## Command Reference

### Global Commands
```
csv2parquet                     # Root command
  ├── parquet <input> <output>  # Convert CSV to Parquet
  └── csv <input> <output>      # Convert Parquet to CSV
```

### Available Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--compression` | `-c` | int | 0 | Compression type (0=UNCOMPRESSED, 1=SNAPPY, 2=GZIP, 3=LZO) |
| `--delimiter` | `-d` | string | "," | Field delimiter for CSV files |
| `--flush` | `-f` | int | 10000 | Number of rows to process before flushing to disk |
| `--verbose` | `-v` | bool | false | Show detailed statistics and performance metrics |
| `--help` | `-h` | bool | false | Display help information |

### Help Commands
```bash
./csv2parquet --help                   # General help
./csv2parquet parquet --help           # CSV to Parquet help
./csv2parquet csv --help               # Parquet to CSV help
```

## Examples

### Basic Conversion Examples
```bash
# CSV to Parquet with default settings
./csv2parquet parquet data.csv

# Parquet to CSV with custom delimiter
./csv2parquet csv data.parquet --delimiter ";"

# CSV to Parquet with compression and verbose output
./csv2parquet parquet large_dataset.csv --compression 1 --verbose
```

### Advanced Usage
```bash
# Process large files with custom flush interval
./csv2parquet parquet big_file.csv big_file.parquet \
  --flush 50000 \
  --compression 2 \
  --verbose

# Convert with pipe delimiter and detailed stats
./csv2parquet csv analytics.parquet analytics.csv \
  --delimiter "|" \
  --flush 1000 \
  --verbose
```

## Performance Features

- **Batch Processing**: Configurable row batch sizes for optimal memory usage
- **Compression**: Support for multiple compression algorithms (SNAPPY, GZIP, LZO)
- **Memory Management**: Efficient memory pooling and garbage collection
- **Progress Tracking**: Runtime statistics including processing time and memory usage
- **Schema Optimization**: Automatic type inference and schema generation

## File Format Support

### CSV Features
- Custom delimiters (comma, semicolon, pipe, tab, etc.)
- Header row detection and processing
- Automatic type inference
- Large file handling with streaming

### Parquet Features
- Columnar storage optimization
- Schema preservation
- Multiple compression algorithms
- Efficient read/write operations
- Row group size optimization (128MB default)

## Development

### Project Structure
```
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command definition
│   ├── csv2parquet.go     # CSV to Parquet conversion
│   └── parquet2csv.go     # Parquet to CSV conversion
├── internal/
│   ├── file/              # File operations and I/O
│   ├── helper/            # Utility functions
│   └── schema/            # Schema management
└── main.go                # Application entry point
```

### Running Tests
```bash
go test ./...                 # Run all tests
go test -v ./...             # Verbose test output
go test -bench . ./...       # Run benchmarks
```

### Building
```bash
go build -o csv2parquet main.go   # Build binary
make build                        # Using Makefile (if available)
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [xitongsys/parquet-go](https://github.com/xitongsys/parquet-go) for Parquet file handling
- CLI powered by [spf13/cobra](https://github.com/spf13/cobra)
- High-performance JSON processing with [bytedance/sonic](https://github.com/bytedance/sonic)