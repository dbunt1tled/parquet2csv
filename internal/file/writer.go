package file

import (
	"encoding/csv"
	"errors"
	"os"
)

type CSVWriter struct {
	file      *os.File
	writer    *csv.Writer
	delimiter string
	flush     int
	idx       int
}

func NewCSVWriter(path string, delimiter string, flush int) (*CSVWriter, error) {
	comma := []rune(delimiter)
	if len(comma) != 1 {
		return nil, errors.New("delimiter must be a single character")
	}
	if flush < 1 {
		return nil, errors.New("flush must be at least 1")
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(f)
	w.Comma = comma[0]
	return &CSVWriter{
		file:      f,
		writer:    w,
		delimiter: delimiter,
		idx:       0,
		flush:     flush,
	}, nil
}

func (w *CSVWriter) WriteS(row []string) error {
	if err := w.writer.Write(row); err != nil {
		return err
	}
	w.idx++
	if (w.idx % w.flush) == 0 {
		w.writer.Flush()
	}
	return nil
}

func (w *CSVWriter) Close() error {
	w.writer.Flush()
	// The flush error must not skip the close, or a failed flush leaks the descriptor.
	return errors.Join(w.writer.Error(), w.file.Close())
}
