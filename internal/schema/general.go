package schema

import (
	"github.com/bytedance/sonic"
	"github.com/iancoleman/strcase"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

type Processor func(record []string, sc interface{}, header []string) interface{}

func ProcessDefault(header []string) (interface{}, Processor) {
	sc := MakeDefaultSchema(header)
	return sc, func(record []string, sc interface{}, header []string) interface{} {
		data := make(map[string]interface{})
		if len(header) != len(record) {
			panic("header and record length not equal")
		}
		for i := range header {
			data[header[i]] = record[i]
		}
		jsonString, _ := sonic.ConfigFastest.Marshal(data)

		err := sonic.ConfigFastest.Unmarshal(jsonString, &sc)
		if err != nil {
			panic(err)
		}
		return sc
	}
}

func MakeDefaultSchema(header []string) interface{} {
	sc := dynamicstruct.NewStruct()
	for i := range header {
		sc.AddField(
			strcase.ToCamel(header[i]),
			"",
			`json:"`+header[i]+`" parquet:"name=`+header[i]+`, type=BYTE_ARRAY, convertedtype=UTF8"`,
		)
	}
	return sc.Build().New()
}
