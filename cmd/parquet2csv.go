package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dbunt1tled/parquet2csv/internal/file"
	"github.com/dbunt1tled/parquet2csv/internal/helper"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/source"
)

var parquet2csv = &cobra.Command{ //nolint:gochecknoglobals // need for init command
	Use:   "csv <input> <output>",
	Short: "Convert parquet to csv",
	Long:  "Convert file from parquet to csv",
	Args:  cobra.RangeArgs(1, 2), //nolint:mnd // args count
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err                error
			input, output, ext string
			delimiter, column  string
			flush              int
			verbose            bool
			fw                 *file.CSVWriter
			fr                 source.ParquetFile
			pr                 *reader.ParquetReader
			rows               []interface{}
			columns, record    []string
			m                  map[string]interface{}
			ok                 bool
		)
		startTime := time.Now()

		input = args[0]
		ext = filepath.Ext(input)
		if ext != ".parquet" {
			return errors.New("file is not parquet file")
		}
		if _, err = file.IsExist(input); err != nil {
			return errors.Wrap(err, "input file "+input+" not exist")
		}
		output = strings.TrimSuffix(input, ext) + ".csv"

		if len(args) == 2 { //nolint:mnd // args count
			output = args[1]
			output = strings.TrimSuffix(output, ".csv") + ".csv"
		}

		flush, err = cmd.Flags().GetInt("flush")
		if err != nil {
			return errors.Wrap(err, "error read flush")
		}
		delimiter, err = cmd.Flags().GetString("delimiter")
		if err != nil {
			return errors.Wrap(err, "error read delimiter")
		}
		verbose, err = cmd.Flags().GetBool("verbose")
		if err != nil {
			return errors.Wrap(err, "error read verbose")
		}

		if _, err = file.IsWritable(filepath.Dir(output)); err != nil {
			return err
		}

		fr, err = local.NewLocalFileReader(input)

		if err != nil {
			return errors.Wrap(err, "error open file reader")
		}

		pr, err = reader.NewParquetReader(fr, nil, 2) //nolint:mnd // nil = generic interface, 2 = goroutines
		if err != nil {
			return errors.Wrap(err, "error open parquet reader")
		}
		defer pr.ReadStop()

		num := int(pr.GetNumRows())
		if num == 0 {
			_, err = file.Create(output)
			if err != nil {
				return errors.Wrap(err, "error create file")
			}
			return nil
		}

		fw, err = file.NewCSVWriter(output, delimiter, flush)
		if err != nil {
			return errors.Wrap(err, "error open file writer")
		}
		defer func(fw *file.CSVWriter) {
			err = fw.Close()
			if err != nil {
				err = errors.Wrap(err, "error close file writer")
			}
		}(fw)
		if err != nil {
			return err
		}
		header := pr.SchemaHandler.SchemaElements
		stringPool := sync.Pool{
			New: func() interface{} {
				s := make([]string, 0, len(header))
				return &s
			},
		}
		recordPtr, ok := stringPool.Get().(*[]string)
		if !ok || recordPtr == nil {
			return errors.New("unexpected record type")
		}
		record = (*recordPtr)[:0]
		for _, el := range header {
			if el.NumChildren != nil {
				continue
			}
			column = el.GetName()
			columns = append(columns, column)
			record = append(record, strings.ToLower(column))
		}
		err = fw.WriteS(record)
		if err != nil {
			return errors.Wrap(err, "error write header")
		}
		stringPool.Put(&record)

		readRows := 0
		for readRows < num {
			rowsToRead := flush
			if num-readRows < flush {
				rowsToRead = num - readRows
			}

			rows, err = pr.ReadByNumber(rowsToRead)
			if err != nil {
				return errors.Wrap(err, "error read rows")
			}
			for _, row := range rows {
				m, err = helper.StructToMap(row)
				if err != nil {
					return fmt.Errorf("unexpected row type: %T, expected map[string]interface{}", row)
				}
				recordPtr, ok = stringPool.Get().(*[]string)
				if !ok || recordPtr == nil {
					return errors.New("unexpected record type")
				}
				record = *recordPtr
				record = record[:0]

				for _, col := range columns {
					record = append(record, helper.AnyToString(m[col]))
				}
				err = fw.WriteS(record)
				if err != nil {
					return errors.Wrap(err, "error write row")
				}
				stringPool.Put(&record)
			}

			readRows += rowsToRead
		}
		if verbose {
			fmt.Printf("%s\n", helper.RuntimeStatistics(startTime, input)) //nolint:forbidigo // verbose output
		}
		return nil
	},
}

//nolint:gochecknoinits // need for init command
func init() {
	rootCmd.AddCommand(parquet2csv)
	parquet2csv.Flags().IntP("compression", "c", 0, "Type of compression")
	parquet2csv.Flags().IntP("flush", "f", file.FlushCount, "number of rows to flush")
	parquet2csv.Flags().StringP("delimiter", "d", ",", "Delimiter for csv file")
	parquet2csv.Flags().BoolP("verbose", "v", false, "Show debug information")
}
