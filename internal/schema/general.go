package schema

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/iancoleman/strcase"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

type Processor func(record []string, sc interface{}, header []string, dataPool *sync.Pool) interface{}

func ProcessDefault(header []string) (interface{}, Processor, error) {
	sc, err := MakeDefaultSchema(header)
	if err != nil {
		return nil, nil, err
	}
	return sc, func(record []string, sc interface{}, header []string, dataPool *sync.Pool) interface{} {
		var (
			dataPtr *map[string]interface{}
			data    map[string]interface{}
			ok      bool
		)

		if dataPool != nil {
			dataPtr, ok = dataPool.Get().(*map[string]interface{})
			if !ok || dataPtr == nil {
				panic("unexpected data type")
			}
			data = *dataPtr
			clear(data)
		} else {
			data = make(map[string]interface{})
		}

		if len(header) != len(record) {
			panic("header and record length not equal")
		}
		for i := range header {
			data[header[i]] = record[i]
		}
		jsonString, _ := sonic.ConfigFastest.Marshal(data)
		if fillErr := sonic.ConfigFastest.Unmarshal(jsonString, &sc); fillErr != nil {
			panic(fillErr)
		}
		if dataPool != nil {
			*dataPtr = data
			dataPool.Put(dataPtr)
		}
		return sc
	}, nil
}

func MakeDefaultSchema(header []string) (interface{}, error) {
	if err := validateHeader(header); err != nil {
		return nil, err
	}

	sc := dynamicstruct.NewStruct()
	used := make(map[string]bool, len(header))
	for i := range header {
		name := fieldName(header[i], i)
		for used[name] {
			name += strconv.Itoa(i + 1)
		}
		used[name] = true

		sc.AddField(
			name,
			"",
			`json:"`+header[i]+`" parquet:"name=`+header[i]+`, type=BYTE_ARRAY, convertedtype=UTF8"`,
		)
	}
	return sc.Build().New(), nil
}

// validateHeader rejects column names that cannot be carried through unchanged. The tags are
// built by concatenation, so a quote silently drops the column and a comma breaks tag parsing,
// and duplicates collapse into one value because the record is keyed by column name.
func validateHeader(header []string) error {
	seen := make(map[string]int, len(header))
	for i := range header {
		column := header[i]
		if column == "" {
			return fmt.Errorf("column %d has an empty name", i+1)
		}
		if strings.ContainsAny(column, `",`) {
			return fmt.Errorf(
				`column %d (%q) contains a quote or a comma, which a parquet tag cannot carry`,
				i+1, column,
			)
		}
		if first, ok := seen[column]; ok {
			return fmt.Errorf("column %d duplicates column %d (%q)", i+1, first+1, column)
		}
		seen[column] = i
	}
	return nil
}

// fieldName derives a valid exported Go identifier for the in-memory struct. Names that carry
// no ASCII letters (Cyrillic, CJK) or start with a digit get a positional one instead. The
// column keeps its own name in the json and parquet tags, so this never reaches the file.
func fieldName(column string, index int) string {
	var b strings.Builder
	for _, r := range strcase.ToCamel(column) {
		if isIdentRune(r) {
			b.WriteRune(r)
		}
	}

	name := b.String()
	if name == "" || !isLetter(rune(name[0])) {
		name = "Col" + strconv.Itoa(index+1) + name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func isIdentRune(r rune) bool {
	return isLetter(r) || (r >= '0' && r <= '9') || r == '_'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
