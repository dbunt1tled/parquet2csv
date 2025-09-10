package helper

import (
	"fmt"
	"os"
	"reflect"
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
	fInfo, _ := os.Stat(inputFile)
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

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", val.Kind())
	}

	result := make(map[string]interface{})
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		result[field.Name] = fieldValue.Interface()
	}

	return result, nil
}

func AnyToString(a any) string {
	switch value := a.(type) {
	case nil:
		return ""
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	case int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(value).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(value).Uint(), 10)
	case []byte:
		return string(value)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case time.Time:
		return value.Format(time.RFC3339)
	default:
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			elem := reflect.ValueOf(value).Elem()
			if !elem.IsValid() {
				return ""
			}
			return AnyToString(elem.Interface())
		}
		return fmt.Sprintf("%v", value)
	}
}
