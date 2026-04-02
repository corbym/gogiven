package model

import (
	"github.com/corbym/gogiven/generator"
)

type testResults struct {
	ID         string `json:"id"`
	Failed     bool   `json:"failed"`
	Skipped    bool   `json:"skipped"`
	TestOutput string `json:"test_output"`
}

func newTestResults(data generator.TestResult) testResults {
	return testResults{
		Failed:     data.Failed,
		TestOutput: data.TestOutput,
		ID:         data.TestID,
		Skipped:    data.Skipped,
	}
}
