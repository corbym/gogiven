package jsonspec_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/corbym/gocrest"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/generator"
	"github.com/corbym/gogiven/generator/jsonspec"
	"testing"
)

var jsonString string
var underTest generator.GoGivensOutputGenerator

func init() {
	underTest = jsonspec.NewTestOutputGenerator()
}

func TestTestOutputGenerator_Generate(testing *testing.T) {
	fileIsConverted()
	AssertThat(testing, jsonString, isValidJSON())
	AssertThat(testing, jsonString, is.EqualToIgnoringWhitespace(
		`{
			"title": "Generator Test",
			"test_state": [
				{
					"test_results": {
						"id": "abc2124",
						"failed": true,
						"skipped": true,
						"test_output": "well alrighty then"
					},
					"test_title": "test title",
					"interesting_givens": {
						"faff": "flap"
					},
					"captured_io": {
						"foob": "barb"
					},
					"given_when_then": {
						"given_when_then": [
							"given",
							"when",
							"then"
						],
						"comment": [
							"Fooing is best",
							"done with friends"
						]
					}
				}
			]
		}`))
}

func isValidJSON() *gocrest.Matcher {
	matcher := &gocrest.Matcher{Describe: "valid jsonString"}
	matcher.Matches = func(actual interface{}) bool {
		buffer := &bytes.Buffer{}
		buffer.WriteString(actual.(string))
		var f interface{}
		err := json.Unmarshal(buffer.Bytes(), &f)
		return err == nil
	}
	return matcher
}

func TestTestOutputGenerator_GenerateConcurrently(testing *testing.T) {
	data := newPageData(false, false)
	for i := 0; i < 15; i++ {
		go func() {
			jsonContent := underTest.Generate(data)
			buffer := new(bytes.Buffer)
			buffer.ReadFrom(jsonContent)
			AssertThat(testing, buffer.String(), is.ValueContaining("Generator Test"))
		}()
	}
}

func TestTestOutputGenerator_FileExtension(t *testing.T) {
	AssertThat(t, underTest.ContentType(), is.EqualTo("application/json"))
}

func TestTestOutputGenerator_Errors(t *testing.T) {
	localUnderTest := jsonspec.NewTestOutputGenerator()

	jsonMarshaller := localUnderTest.MarshalJSON
	defer func() {
		recovered := recover()
		localUnderTest.MarshalJSON = jsonMarshaller
		AssertThat(t, recovered, is.Not(is.Nil()))
	}()
	localUnderTest.MarshalJSON = func(v interface{}) ([]byte, error) {
		return nil, errors.New("bugger")
	}

	localUnderTest.Generate(newPageData(false, true))
}

func fileIsConverted() {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(underTest.Generate(newPageData(true, true)))
	jsonString = buffer.String()
}

func newPageData(skipped bool, failed bool) generator.PageData {
	var testData []generator.TestData

	capturedIO := make(map[interface{}]interface{})
	capturedIO["foob"] = "barb"
	interestingGivens := make(map[interface{}]interface{})
	interestingGivens["faff"] = "flap"
	testData = append(testData, generator.TestData{
		TestTitle: "test title",
		ParsedTestContent: base.ParsedTestContent{
			GivenWhenThen: []string{"given", "when", "then"},
			Comment:       []string{"Fooing is best", "done with friends"},
		},
		CapturedIO:        capturedIO,
		InterestingGivens: interestingGivens,
		TestResult: generator.TestResult{
			Failed:     failed,
			Skipped:    skipped,
			TestOutput: "well alrighty then",
			TestID:     "abc2124",
		},
	})
	return generator.PageData{
		TestData: testData,
		Title:    "Generator Test",
	}
}
