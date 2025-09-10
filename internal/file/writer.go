package file

import (
	"encoding/csv"
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
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(f)
	w.Comma = rune(delimiter[0])
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
	if err := w.file.Close(); err != nil {
		return err
	}
	return w.file.Close()
}
