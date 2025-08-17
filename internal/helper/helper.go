package helper

import (
	"fmt"
	"github.com/dbunt1tled/parquet2csv/internal/file"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func StrToInt64(str string, panicIfErr bool) int64 {
	if str == "" {
		return 0
	}
	i, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		if panicIfErr {
			panic(err)
		}
		return 0
	}
	return i
}

func StrToInt32(str string, panicIfErr bool) int32 {
	if str == "" {
		return 0
	}
	i, err := strconv.ParseInt(strings.TrimSpace(str), 10, 32)
	if err != nil {
		if panicIfErr {
			panic(err)
		}
		return 0
	}
	return int32(i)
}

func ConvertToFloat(str string, panicIfErr bool) float64 {
	if str == "" {
		return 0
	}
	i, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		if panicIfErr {
			panic(err)
		}
		return 0
	}
	return i
}

func RuntimeStatistics(startTime time.Time, inputFile string) string {
	pc, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	fInfo, _ := file.Info(inputFile)
	return fmt.Sprintf(
		"%s (%s): %s Processed %s (%s)",
		inputFile,
		GetFileSize(fInfo.Size()),
		name,
		time.Since(startTime).Round(time.Second).String(),
		MemoryUsage(),
	)
}

func MemoryUsage() string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return fmt.Sprintf(
		"TotalAlloc: %v MB, Sys: %v MB",
		memStats.TotalAlloc/1024/1024, //nolint:mnd // Convert to MB
		memStats.Sys/1024/1024,        //nolint:mnd // Convert to MB
	)
}

func AppHelp(help bool) {
	if help {
		log.Printf("%s", `
https://github.com/dbunt1tled/parquet2csv
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
`)
		os.Exit(0)
	}
}

func GetFileSize(size int64) string {
	sizeKb := 1024.0
	sizeMb := sizeKb * sizeKb
	sizeGb := sizeMb * sizeKb

	switch {
	case float64(size) < sizeMb:
		return fmt.Sprintf("%.2f Kb", float64(size)/sizeKb)
	case float64(size) < sizeGb:
		return fmt.Sprintf("%.2f Mb", float64(size)/sizeMb)
	default:
		return fmt.Sprintf("%.2f Gb", float64(size)/sizeGb)
	}
}
