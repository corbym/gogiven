package gogiven

import (
	"fmt"
)

type TestMetaData struct {
	TestId     string
	Failed     bool
	TestOutput string
}

func newTestMetaData(testName string) *TestMetaData {
	testContext := new(TestMetaData)
	testContext.TestId = testName
	testContext.Failed = false
	return testContext
}
func (t *TestMetaData) Logf(format string, args ...interface{}) {
	t.TestOutput = fmt.Sprintf(format, args...)
	t.Failed = true
}

func (t *TestMetaData) Errorf(format string, args ...interface{}) {
	t.TestOutput = fmt.Sprintf(format, args...)
	t.Failed = true
}

func (t *TestMetaData) FailNow() {
	t.Failed = true
}

func (t *TestMetaData) Helper() {
	// do nothing
}
