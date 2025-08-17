# csv2parquet

[![Go Version](https://img.shields.io/badge/go-1.21+-blue?logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/dbunt1tled/parquet2csv.svg)](https://pkg.go.dev/github.com/dbunt1tled/parquet2csv)
[![Build Status](https://github.com/dbunt1tled/parquet2csv/workflows/Build/badge.svg)](https://github.com/dbunt1tled/parquet2csv/actions)
[![Release](https://img.shields.io/github/v/release/dbunt1tled/parquet2csv)](https://github.com/dbunt1tled/github-unfollow/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/dbunt1tled/parquet2csv)](https://goreportcard.com/report/github.com/dbunt1tled/parquet2csv)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example) ![example workflow](https://github.com/dbunt1tled/parquet2csv/actions/workflows/build.yml/badge.svg)

A fast, reliable CLI tool for converting CSV files to Apache Parquet format. Built in Go, it’s designed for data workflows that need efficient, schema-aware columnar storage.

## Install & Run

```
$ go build
$ ./csv2parquet file.csv file.parquet
```

## Additional Flags

```
$ ./csv2parquet --help
Usage of ./csv2parquet:
  --delimiter string
        Delimiter for csv file (default ",")
  --flush int
        number of rows to flush (default 10000)
  --help
        Show this help message
  --schema string
        schema of csv file (default "default")
  --compression int
        Type of compression (default 0)
  --verbose
        Statistic info in the end
  <input file path>
  <output file path>
```