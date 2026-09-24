package schema

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

// Processor fills the schema struct from one CSV record and returns it. The same pointer comes
// back every time, so the record must be handed to the writer before the next one is processed.
type Processor func(record []string) interface{}

func ProcessDefault(header []string) (interface{}, Processor, error) {
	sc, err := MakeDefaultSchema(header)
	if err != nil {
		return nil, nil, err
	}

	// MakeDefaultSchema adds one field per column in header order, so column i is field i and
	// the record can be written straight into the struct without a name lookup.
	fields := reflect.ValueOf(sc).Elem()
	columns := len(header)

	return sc, func(record []string) interface{} {
		if len(record) != columns {
			panic("header and record length not equal")
		}
		for i := range record {
			fields.Field(i).SetString(record[i])
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
