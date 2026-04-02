package jsonspec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/corbym/gogiven/generator"
	"github.com/corbym/gogiven/generator/jsonspec/model"
	"io"
)

// TestOutputGenerator is an implementation of the GoGivensOutputGenerator that generates a JSON file per
// test. It is thread safe between goroutines.
type TestOutputGenerator struct {
	generator.GoGivensOutputGenerator
	MarshalJSON func(v interface{}) ([]byte, error)
}

// NewTestOutputGenerator creates a template that is used to generate the json output.
func NewTestOutputGenerator() *TestOutputGenerator {
	outputGenerator := new(TestOutputGenerator)
	outputGenerator.MarshalJSON = json.Marshal
	return outputGenerator
}

// ContentType for the output generated.
func (outputGenerator *TestOutputGenerator) ContentType() string {
	return "application/json"
}

// Generate generates json output for a test. The return string contains the json
// that goes into the output file generated in gogivens.GenerateTestOutput().
// The function panics if the template cannot be generated.
func (outputGenerator *TestOutputGenerator) Generate(pageData generator.PageData) io.Reader {
	jsonBytes, err := outputGenerator.MarshalJSON(model.NewJSONData(pageData))
	if err != nil {
		panic(fmt.Sprintf("Could not marshal pageData to json: %v", err))
	}
	var out = new(bytes.Buffer)
	json.Indent(out, jsonBytes, "", "\t")
	return out
}

// GenerateIndex generates an index of all the data from the tests in JSON format.
func (outputGenerator *TestOutputGenerator) GenerateIndex(indexData []generator.IndexData) io.Reader {
	return bytes.NewReader(nil) //TODO: implement some kind of index
}
