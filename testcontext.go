package gogiven

import (
	"fmt"
)
// TestMetaData holds some information about the test, it's id, whether it failed or not
// and the output it sent to t.Errorf or t.Error etc.
type TestMetaData struct {
	TestId     string
	Failed     bool
	TestOutput string
}
//NewTestMetaData creates a new TestMetaData object. Used internally.
func NewTestMetaData(testName string) *TestMetaData {
	testContext := new(TestMetaData)
	testContext.TestId = testName
	testContext.Failed = false
	return testContext
}
//Logf marks this test as failed and sets the test output to the formatted string. 
func (t *TestMetaData) Logf(format string, args ...interface{}) {
	t.TestOutput = fmt.Sprintf(format, args...)
	t.Failed = true
}
//Errorf marks this test as failed and sets the test output to the formatted string. 
func (t *TestMetaData) Errorf(format string, args ...interface{}) {
	t.TestOutput = fmt.Sprintf(format, args...)
	t.Failed = true
}
//FailNow marks this test as failed. 
func (t *TestMetaData) FailNow() {
	t.Failed = true
}
//Helper does nothing. It's just in case some package that consumes t 
// calls it.
func (t *TestMetaData) Helper() {
	// do nothing
}
func (t *TestMetaData) Name() string {
	return t.TestId
}
