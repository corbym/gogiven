package base

import (
	"github.com/corbym/gogiven/testdata"
	"sync"
)

// Some holds the test context and has a reference to the test's testing.T
type Some struct {
	sync.RWMutex
	globalTestingT    TestingT
	TestingT          *TestMetaData
	TestTitle         string
	interestingGivens testdata.InterestingGivens
	capturedIO        testdata.CapturedIO
	GivenWhenThen     string
}

//NewSome creates a new Some context. This is an internal function that was exported for testing.
func NewSome(
	globalTestingT TestingT,
	testTitle string,
	testContext *TestMetaData,
	givenWhenThen string,
	givenFunc ...GivenData) *Some {

	some := new(Some)
	some.TestingT = testContext
	some.TestTitle = testTitle
	some.globalTestingT = globalTestingT
	some.GivenWhenThen = givenWhenThen
	some.interestingGivens = make(testdata.InterestingGivens)
	some.capturedIO = make(testdata.CapturedIO)

	if len(givenFunc) > 0 {
		for _, someGivenFunc := range givenFunc {
			someGivenFunc(some.interestingGivens)
		}
	}
	return some
}

// CapturedIO is a convenience method for retrieving the CapturedIO map
func (some *Some) CapturedIO() map[string]interface{} {
	some.Lock()
	defer some.Unlock()
	return copyMap(some.capturedIO)
}

// InterestingGivens is a convenience method for retrieving the InterestingGivens map
func (some *Some) InterestingGivens() map[string]interface{} {
	some.Lock()
	defer some.Unlock()
	return copyMap(some.interestingGivens)
}

// When - call When when you want to perform some action, call a function, or perform a test operation.
func (some *Some) When(action ...CapturedIOGivenData) *Some {
	some.Lock()
	defer some.Unlock()
	action[0](some.capturedIO, some.interestingGivens) // TODO: there could be multiple actions..
	return some
}

// Then is a function that executes the given function and asserts whether the test has failed.
// It can be called in a table test (for loop). Provide a function in which assertions will be made.
// Use the TestingT typed var in place of testing.T.
// The test state is recorded in TestingT type and goGiven fails the test if the error methods (ErrorF etc)
// were called after the function exits.
func (some *Some) Then(assertions TestingWithGiven) *Some {
	some.Lock()
	defer some.Unlock()
	assertions(some.TestingT, some.capturedIO, some.interestingGivens)
	if some.TestingT.failed {
		globalTestingT := some.globalTestingT
		globalTestingT.Helper()
		globalTestingT.Errorf(some.TestingT.TestOutput)
	}
	return some
}

func copyMap(ios map[string]interface{}) map[string]interface{} {
	var newMap = make(map[string]interface{})
	for k, v := range ios {
		newMap[k] = v.(interface{})
	}
	return newMap
}
