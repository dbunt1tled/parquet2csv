package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/dbunt1tled/parquet2csv/internal/file"
	"github.com/dbunt1tled/parquet2csv/internal/helper"
	"github.com/dbunt1tled/parquet2csv/internal/schema"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

const (
	kilobyte = 1024
	megabyte = 1024 * kilobyte
)

var csv2parquet = &cobra.Command{ //nolint:gochecknoglobals // need for init command
	Use:   "parquet <input> <output>",
	Short: "Convert csv to parquet",
	Long:  "Convert file from csv to parquet",
	Args:  cobra.RangeArgs(1, 2), //nolint:mnd // args count
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err                error
			input, output, ext string
			compression        int
			delimiter          string
			flush              int
			rowGroupSize       int
			pageSize           int
			verbose            bool
			structType         interface{}
			processor          schema.Processor
			fw                 source.ParquetFile
			pw                 *writer.ParquetWriter
		)
		startTime := time.Now()

		input = args[0]
		ext = filepath.Ext(input)
		if ext != ".csv" {
			return errors.New("file is not csv file")
		}
		if _, err = file.IsExist(input); err != nil {
			return errors.Wrap(err, "input file "+input+" not exist")
		}
		output = strings.TrimSuffix(input, ext) + ".parquet"

		if len(args) == 2 { //nolint:mnd // args count
			output = args[1]
			output = strings.TrimSuffix(output, ".parquet") + ".parquet"
		}

		compression, err = cmd.Flags().GetInt("compression")
		if err != nil {
			return errors.Wrap(err, "error read compression")
		}
		// parquet-go silently compresses to nothing for a codec it has no compressor for,
		// which writes a file with a schema and no rows, so only registered codecs are allowed.
		supported := []int{
			int(parquet.CompressionCodec_UNCOMPRESSED),
			int(parquet.CompressionCodec_SNAPPY),
			int(parquet.CompressionCodec_GZIP),
			int(parquet.CompressionCodec_LZ4),
			int(parquet.CompressionCodec_ZSTD),
		}
		if !slices.Contains(supported, compression) {
			return fmt.Errorf("unsupported compression %d, want one of %v", compression, supported)
		}
		flush, err = cmd.Flags().GetInt("flush")
		if err != nil {
			return errors.Wrap(err, "error read flush")
		}
		if flush < 1 {
			return fmt.Errorf("flush must be at least 1, got %d", flush)
		}
		rowGroupSize, err = cmd.Flags().GetInt("row-group-size")
		if err != nil {
			return errors.Wrap(err, "error read row-group-size")
		}
		if rowGroupSize < 1 {
			return fmt.Errorf("row-group-size must be at least 1 MB, got %d", rowGroupSize)
		}
		pageSize, err = cmd.Flags().GetInt("page-size")
		if err != nil {
			return errors.Wrap(err, "error read page-size")
		}
		if pageSize < 1 {
			return fmt.Errorf("page-size must be at least 1 KB, got %d", pageSize)
		}
		delimiter, err = cmd.Flags().GetString("delimiter")
		if err != nil {
			return errors.Wrap(err, "error read delimiter")
		}
		if len([]rune(delimiter)) != 1 {
			return fmt.Errorf("delimiter must be a single character, got %q", delimiter)
		}
		verbose, err = cmd.Flags().GetBool("verbose")
		if err != nil {
			return errors.Wrap(err, "error read verbose")
		}

		if _, err = file.IsWritable(filepath.Dir(output)); err != nil {
			return err
		}

		reader, err := file.NewRowReader(input, []rune(delimiter)[0], false)
		if err != nil {
			return err
		}
		defer reader.Close()

		fw, err = local.NewLocalFileWriter(output)
		if err != nil {
			return err
		}
		// Best effort for the error paths; the success path closes and checks the error below.
		defer fw.Close()

		rows := 0
		for {
			record, readErr := reader.Next()
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {
					break
				}
				return errors.Wrap(readErr, "read error")
			}

			// The first record is the header and there is nothing to write yet, so the writer
			// staying nil is also what marks the header as not seen.
			if pw == nil {
				// The reader hands back the same slice every time, so the header has to be
				// copied before the next record overwrites it.
				structType, processor, err = schema.ProcessDefault(slices.Clone(record))
				if err != nil {
					return errors.Wrap(err, "build schema")
				}
				pw, err = writer.NewParquetWriter(fw, structType, 2) //nolint:mnd // maybe the number of threads
				if err != nil {
					return errors.Wrap(err, "can't create parquet writer")
				}
				pw.RowGroupSize = int64(rowGroupSize) * megabyte
				// One ColumnIndex and OffsetIndex entry is kept in memory per page until
				// WriteStop, so the page size — not the row group size — is what decides
				// whether peak memory grows with the length of the input.
				pw.PageSize = int64(pageSize) * kilobyte
				pw.CompressionType = parquet.CompressionCodec(int32(compression))
				continue
			}

			if err = pw.Write(processor(record)); err != nil {
				return errors.Wrap(err, "write error")
			}

			rows++
			if rows%flush == 0 {
				// Flush(false) only turns the buffered rows into pages; it closes a row group
				// once RowGroupSize is reached and not before. Forcing a row group here instead
				// would keep per-row-group footer metadata for every flush until WriteStop.
				if err = pw.Flush(false); err != nil {
					return errors.Wrap(err, "write flush error")
				}
			}
		}

		if pw == nil {
			return errors.New("input file " + input + " is empty: no header row")
		}

		if err = pw.WriteStop(); err != nil {
			return errors.Wrap(err, "write stop error")
		}

		if err = fw.Close(); err != nil {
			return errors.Wrap(err, "close writer error")
		}
		if err = reader.Close(); err != nil {
			return errors.Wrap(err, "close reader error")
		}
		if verbose {
			fmt.Printf("%s\n", helper.RuntimeStatistics(startTime, input)) //nolint:forbidigo  // verbose output
		}
		return nil
	},
}

//nolint:gochecknoinits // need for init command
func init() {
	rootCmd.AddCommand(csv2parquet)
	csv2parquet.Flags().IntP("compression", "c", 0, "Type of compression")
	csv2parquet.Flags().IntP("flush", "f", file.FlushCount, "number of rows to flush")
	csv2parquet.Flags().IntP("row-group-size", "r", 8, "Target row group size in MB") //nolint:mnd // default 8MB
	csv2parquet.Flags().IntP("page-size", "p", 1024, "Target data page size in KB")   //nolint:mnd // default 1MB
	csv2parquet.Flags().StringP("delimiter", "d", ",", "Delimiter for csv file")
	csv2parquet.Flags().BoolP("verbose", "v", false, "Show debug information")
}
