package base

import (
	"fmt"
	"sync"
)

// TestMetaData holds some information about the test, it's id, whether it failed or not
// and the output it sent to t.Errorf or t.Error etc.
type TestMetaData struct {
	sync.RWMutex
	TestId     string
	failed     bool
	skipped    bool
	testOutput string
}

//NewTestMetaData creates a new TestMetaData object. Used internally.
func NewTestMetaData(testName string) *TestMetaData {
	testContext := new(TestMetaData)
	testContext.TestId = testName
	testContext.failed = false
	testContext.skipped = false
	return testContext
}

//Logf marks this test as failed and sets the test output to the formatted string.
func (t *TestMetaData) Logf(format string, args ...interface{}) {
	t.Lock()
	defer t.Unlock()
	t.testOutput = fmt.Sprintf(format, args...)
	t.failed = true
}

//Errorf marks this test as failed and sets the test output to the formatted string.
func (t *TestMetaData) Errorf(format string, args ...interface{}) {
	t.Lock()
	defer t.Unlock()
	t.testOutput = fmt.Sprintf(format, args...)
	t.failed = true
}

//FailNow marks this test as failed.
func (t *TestMetaData) FailNow() {
	t.Lock()
	defer t.Unlock()
	t.failed = true
}

//Helper does nothing. It's just in case some package that consumes t
// calls it.
func (t *TestMetaData) Helper() {
	// do nothing
}

//Name returns the id (the test name, possibly with some uniqueid appended)
func (t *TestMetaData) Name() string {
	t.RLock()
	defer t.RUnlock()
	return t.TestId
}

//Failed reports the test has failed to the meta data.
func (t *TestMetaData) Failed() bool {
	t.RLock()
	defer t.RUnlock()
	return t.failed
}

func (t *TestMetaData) Skipf(format string, args ...interface{}) {
	t.Lock()
	defer t.Unlock()
	t.skipped = true
	t.testOutput = fmt.Sprintf(format, args...)
}

func (t *TestMetaData) Skipped() bool {
	t.RLock()
	defer t.RUnlock()
	return t.skipped
}

func (t *TestMetaData) TestOutput() string{
	t.RLock()
	defer t.RUnlock()
	return t.testOutput
}