package base

import (
	"github.com/corbym/gogiven/testdata"
	"sync"
)

// Some holds the test context and has a reference to the test's testing.T
type Some struct {
	sync.RWMutex
	globalTestingT    TestingT
	testMetaData      *TestMetaData
	testTitle         string
	interestingGivens testdata.InterestingGivens
	capturedIO        testdata.CapturedIO
	parsedTestContent ParsedTestContent
}

// NewSome creates a new Some context. This is an internal function that was exported for testing.
func NewSome(
	globalTestingT TestingT,
	testTitle string,
	testingT *TestMetaData,
	givenWhenThen ParsedTestContent,
	givenFunc ...GivenData) *Some {

	some := new(Some)
	some.testMetaData = testingT
	some.testTitle = testTitle
	some.globalTestingT = globalTestingT
	some.parsedTestContent = givenWhenThen
	some.interestingGivens = make(testdata.InterestingGivens)
	some.capturedIO = make(testdata.CapturedIO)

	if len(givenFunc) > 0 {
		for _, someGivenFunc := range givenFunc {
			someGivenFunc(some.interestingGivens)
		}
	}
	return some
}

// TestTitle is the name of the test
func (some *Some) TestTitle() string {
	some.RLock()
	defer some.RUnlock()
	return some.testTitle
}

// ParsedTestContent holds a parsed test as an array string.
// All lines of the test func will be listed, from the first call to the first Given
// up to the end of the test func and converted to (as close as possible) natural language.
func (some *Some) ParsedTestContent() ParsedTestContent {
	some.RLock()
	defer some.RUnlock()

	return some.parsedTestContent
}

// TestMetaData is an interface that mimics testingT but stores the test state rather than act on it.
// Gogivens will act on the meta data's behalf via globalTestingT (the "real" testing.T for the test).
func (some *Some) TestMetaData() *TestMetaData {
	some.RLock()
	defer some.RUnlock()
	return some.testMetaData
}

// CapturedIO is a convenience method for retrieving the CapturedIO map
func (some *Some) CapturedIO() map[interface{}]interface{} {
	some.RLock()
	defer some.RUnlock()
	return some.capturedIO
}

// InterestingGivens is a convenience method for retrieving the InterestingGivens map
func (some *Some) InterestingGivens() map[interface{}]interface{} {
	some.RLock()
	defer some.RUnlock()
	return some.interestingGivens
}

// When - call When when you want to perform Some action, call a function, or perform a test operation.
func (some *Some) When(action ...CapturedIOGivenData) *Some {
	some.Lock()
	defer some.Unlock()
	action[0](some.capturedIO, some.interestingGivens) // TODO: there could be multiple actions..
	return some
}

// Then is a function that executes the given function and asserts whether the test has failed.
// It can be called in a table test (for loop). Provide a function in which assertions will be made.
// Use the testMetaData typed var in place of testing.T.
// The test state is recorded in testMetaData type and goGiven fails the test if the error methods (ErrorF etc)
// were called after the function exits.
func (some *Some) Then(assertions TestingWithGiven) *Some {
	some.Lock()
	defer some.Unlock()
	if !some.testMetaData.Skipped() {
		assertions(some.testMetaData, some.capturedIO, some.interestingGivens)
		if some.testMetaData.Failed() {
			globalTestingT := some.globalTestingT
			globalTestingT.Helper()
			globalTestingT.Errorf(some.testMetaData.TestOutput())
		}
	}
	return some
}

// SkippingThisOne still records we have a skipped tests in our test output generator
func (some *Some) SkippingThisOne(reason string, args ...interface{}) *Some {
	some.testMetaData.Skipf(reason, args...)
	some.globalTestingT.Helper()
	some.globalTestingT.Skipf(reason, args...) // skip so we don't worry about it
	return some
}

// SkippingThisOneIf skips if the condition is true, and still records we have a skipped tests in our test output generator.
// This will be best used in a table test (range) when running sub-tests, since in a main test the entire test will be skipped
// and the condition pointless.
func (some *Some) SkippingThisOneIf(why func(someData ...interface{}) bool, reason string, args ...interface{}) *Some {
	if why() {
		some.testMetaData.Skipf(reason, args...)
		some.globalTestingT.Helper()
		some.globalTestingT.Skipf(reason, args...) // skip so we don't worry about it
	}
	return some
}
