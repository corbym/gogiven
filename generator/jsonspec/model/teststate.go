package model

import (
	"fmt"
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/generator"
)

type parsedTestContent struct {
	GivenWhenThen []string `json:"given_when_then"`
	// Comment contains each comment line from the tests' comment block
	Comment []string `json:"comment"`
}

type testState struct {
	TestResults       testResults            `json:"test_results"`
	TestTitle         string                 `json:"test_title"`
	InterestingGivens map[string]interface{} `json:"interesting_givens"`
	CapturedIO        map[string]interface{} `json:"captured_io"`
	GivenWhenThen     parsedTestContent      `json:"given_when_then"`
}

func newTestState(testData generator.TestData) *testState {
	return &testState{
		TestResults:       newTestResults(testData.TestResult),
		TestTitle:         testData.TestTitle,
		InterestingGivens: convertToMapStringInterface(testData.InterestingGivens),
		CapturedIO:        convertToMapStringInterface(testData.CapturedIO),
		GivenWhenThen:     newParsedTestContent(testData.ParsedTestContent),
	}
}

func newParsedTestContent(content base.ParsedTestContent) parsedTestContent {
	return parsedTestContent{
		GivenWhenThen: content.GivenWhenThen,
		Comment:       content.Comment,
	}
}

func convertToMapStringInterface(someMap map[interface{}]interface{}) (converted map[string]interface{}) {
	converted = make(map[string]interface{}, len(someMap))
	for k, v := range someMap {
		newKey := fmt.Sprintf("%v", k)
		converted[newKey] = v
	}
	return
}
