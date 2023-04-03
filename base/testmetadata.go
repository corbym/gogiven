package base

import (
	"fmt"
	"sync"
)

// TestMetaData holds Some information about the test, it's id, whether it failed or not
// and the output it sent to t.Errorf or t.Error etc.
type TestMetaData struct {
	sync.RWMutex
	TestID     string
	failed     bool
	skipped    bool
	testOutput string
}

// NewTestMetaData creates a new TestMetaData object. Used internally.
func NewTestMetaData(testName string) *TestMetaData {
	testContext := &TestMetaData{
		TestID:  testName,
		failed:  false,
		skipped: false,
	}
	return testContext
}

// Logf marks this test as failed and sets the test output to the formatted string.
func (t *TestMetaData) Logf(format string, args ...interface{}) {
	t.Lock()
	defer t.Unlock()
	t.testOutput = fmt.Sprintf(format, args...)
	t.failed = true
}

// Errorf marks this test as failed and sets the test output to the formatted string.
func (t *TestMetaData) Errorf(format string, args ...interface{}) {
	t.Lock()
	defer t.Unlock()
	t.testOutput = fmt.Sprintf(format, args...)
	t.failed = true
}

// FailNow marks this test as failed.
func (t *TestMetaData) FailNow() {
	t.Lock()
	defer t.Unlock()
	t.failed = true
}

// Fail marks this test as failed.
func (t *TestMetaData) Fail() {
	t.Lock()
	defer t.Unlock()
	t.failed = true
}

// Helper does nothing. It's just in case Some package that consumes t
// calls it.
func (t *TestMetaData) Helper() {
	// do nothing
}

// Name returns the id (the test name, possibly with some uniqueid appended)
func (t *TestMetaData) Name() string {
	t.RLock()
	defer t.RUnlock()
	return t.TestID
}

// Failed reports the test has failed to the metadata.
func (t *TestMetaData) Failed() bool {
	t.RLock()
	defer t.RUnlock()
	return t.failed
}

// Skipf reports the test as skipped with a formatted message to the test meta data
func (t *TestMetaData) Skipf(format string, args ...interface{}) {
	t.Lock()
	defer t.Unlock()
	t.skipped = true
	t.testOutput = fmt.Sprintf(format, args...)
}

// Skipped returns the state of the test meta data, if skipped then it's true
func (t *TestMetaData) Skipped() bool {
	t.RLock()
	defer t.RUnlock()
	return t.skipped
}

// TestOutput returns the recorded test output
func (t *TestMetaData) TestOutput() string {
	t.RLock()
	defer t.RUnlock()
	return t.testOutput
}
