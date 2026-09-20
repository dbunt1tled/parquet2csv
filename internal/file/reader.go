package file

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type BatchProcessor struct {
	batchSize  int
	inputFile  string
	skipHeader bool
	delimiter  rune
	batchChan  chan Batch
	resultChan chan []Row
}

type Row struct {
	data   []string
	rowNum int
}

type Batch struct {
	Rows  [][]string
	Start int
	Id    int
}

func NewBatchProcessor(
	inputFile string,
	batchSize int,
	delimiter rune,
	skipHeader bool,
) *BatchProcessor {
	return &BatchProcessor{
		batchSize:  batchSize,
		inputFile:  inputFile,
		delimiter:  delimiter,
		skipHeader: skipHeader,
	}
}

func (bp *BatchProcessor) Reader() (batchChan chan Batch, errorChan chan error) {
	batchChan = make(chan Batch, 2)
	errorChan = make(chan error, 2)
	go func() {
		defer close(batchChan)

		file, err := os.Open(bp.inputFile)
		if err != nil {
			errorChan <- errors.Wrap(err, "error opening file "+bp.inputFile)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				errorChan <- errors.Wrap(err, "error closing file "+bp.inputFile)
				return
			}
		}(file)

		reader := csv.NewReader(file)
		reader.Comma = bp.delimiter
		if bp.skipHeader {
			if _, err := reader.Read(); err != nil {
				errorChan <- errors.Wrap(err, "error reading header")
				return
			}
		}
		batchID := 0
		for {
			batch := make([][]string, 0, bp.batchSize)
			startRow := batchID*bp.batchSize + 1
			for i := 0; i < bp.batchSize; i++ {
				record, err := reader.Read()
				if err != nil {
					if err == io.EOF {
						break
					}
					errorChan <- errors.Wrap(err, "error reading row "+strconv.Itoa(startRow+i))
					return
				}
				batch = append(batch, record)
			}
			if len(batch) == 0 {
				break
			}
			batchChan <- Batch{
				Rows:  batch,
				Start: startRow,
				Id:    batchID,
			}

			batchID++
		}
	}()
	return batchChan, errorChan
}
