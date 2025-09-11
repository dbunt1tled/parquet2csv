package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
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
			verbose            bool
			header             []string
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

		fw, err = local.NewLocalFileWriter(output)
		if err != nil {
			return err
		}
		defer func(fw source.ParquetFile) {
			err = fw.Close()
			if err != nil {
				err = errors.Wrap(err, "close writer error")
			}
		}(fw)
		if err != nil {
			return err
		}
		i := 0
		bp := file.NewBatchProcessor(input, file.FlushCount, []rune(delimiter)[0], false)
		bCh, eCh := bp.Reader()

		dataPool := &sync.Pool{
			New: func() interface{} {
				data := make(map[string]interface{})
				return &data
			},
		}

		for rows := range bCh {
			select {
			case err = <-eCh:
				return errors.Wrap(err, "write error")
			default:
				for _, rec := range rows.Rows {
					if i == 0 {
						header = rec
						structType, processor = schema.ProcessDefault(header)
						pw, err = writer.NewParquetWriter(fw, structType, 2) //nolint:mnd // maybe the number of threads
						if err != nil {
							return errors.Wrap(err, "can't create parquet writer")
						}
						pw.RowGroupSize = 128 * 1024 * 1024 //nolint:mnd // 128MB
						pw.CompressionType = parquet.CompressionCodec(int32(compression))
						i++
						continue
					}

					eData := processor(rec, structType, header, dataPool)
					if err = pw.Write(eData); err != nil {
						return errors.Wrap(err, "write error")
					}

					if i == flush {
						if err = pw.Flush(true); err != nil {
							return errors.Wrap(err, "write flush error")
						}
						i = 0
					}
					i++
				}
			}
		}

		if err = pw.WriteStop(); err != nil {
			return errors.Wrap(err, "write stop error")
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
	csv2parquet.Flags().StringP("delimiter", "d", ",", "Delimiter for csv file")
	csv2parquet.Flags().BoolP("verbose", "v", false, "Show debug information")
}
