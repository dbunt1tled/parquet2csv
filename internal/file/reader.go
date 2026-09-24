package file

import (
	"bufio"
	"encoding/csv"
	stderrors "errors"
	"io"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// readBufferSize keeps the scanner off the default 4KB buffer, so a large input costs
// megabyte-sized reads instead of thousands of small ones.
const readBufferSize = 1 << 20

// RowReader streams a CSV file one record at a time. The slice returned by Next is only valid
// until the next call to Next, which is what lets the underlying reader keep reusing it.
type RowReader struct {
	file   *os.File
	reader *csv.Reader
	row    int
}

func NewRowReader(inputFile string, delimiter rune, skipHeader bool) (*RowReader, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, errors.Wrap(err, "error opening file "+inputFile)
	}

	reader := csv.NewReader(bufio.NewReaderSize(f, readBufferSize))
	reader.Comma = delimiter
	reader.ReuseRecord = true

	r := &RowReader{file: f, reader: reader, row: 0}
	if skipHeader {
		if _, err = reader.Read(); err != nil {
			return nil, stderrors.Join(errors.Wrap(err, "error reading header"), r.Close())
		}
		r.row++
	}
	return r, nil
}

// Next returns the next record, or io.EOF once the file is exhausted. io.EOF is returned
// unwrapped so the caller can end the loop on it without unwrapping a read failure by mistake.
func (r *RowReader) Next() ([]string, error) {
	record, err := r.reader.Read()
	if err != nil {
		if stderrors.Is(err, io.EOF) {
			return nil, io.EOF
		}
		return nil, errors.Wrap(err, "error reading row "+strconv.Itoa(r.row+1))
	}
	r.row++
	return record, nil
}

func (r *RowReader) Close() error {
	return r.file.Close()
}
