package helper

import (
	"fmt"
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

func RuntimeStatistics(startTime time.Time) string {
	pc, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	return fmt.Sprintf(
		"%s Processed %s (%s)",
		name,
		time.Since(startTime).Round(time.Second).String(),
		MemoryUsage(),
	)
}

func MemoryUsage() string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return fmt.Sprintf("TotalAlloc: %v MB, Sys: %v MB", memStats.TotalAlloc/1024/1024, memStats.Sys/1024/1024)

}