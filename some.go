package gogiven

type TestingT interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
	Name() string
}

// Some holds the test context and has a reference to the test's testing.T
type Some struct {
	globalTestingT    TestingT
	testingT          *TestMetaData
	interestingGivens *InterestingGivens
	capturedIO        *CapturedIO
	givenWhenThen	  string
}

//NewSome creates a new Some context. This is an internal function that was exported for testing.
func NewSome(
	globalTestingT TestingT,
	testContext *TestMetaData,
	givenWhenThen string,
	givenFunc ...func(givens *InterestingGivens)) *Some {

	some := new(Some)
	some.testingT = testContext
	some.globalTestingT = globalTestingT
	some.givenWhenThen = givenWhenThen
	some.capturedIO = newCapturedIO()
	givens := newInterestingGivens()

	if len(givenFunc) > 0 {
		for _, someGivenFunc := range givenFunc {
			someGivenFunc(givens)
		}
	}
	some.interestingGivens = givens
	return some
}

// Convenience method for retrieving the CapturedIO map
func (some *Some) CapturedIO() map[string]interface{} {
	return some.capturedIO.CapturedIO
}

// Convenience method for retrieving the InterestingGivens map
func (some *Some) InterestingGivens() map[string]interface{} {
	return some.interestingGivens.Givens
}

// Call When when you want to perform some action, call a function, or perform a test operation.
func (some *Some) When(action ...func(actual *CapturedIO, givens *InterestingGivens)) *Some {
	action[0](some.capturedIO, some.interestingGivens) // TODO: there could be multiple actions..
	return some
}

// Call Then to perform assersions. Provide a function in which assertions will be made.
func (some *Some) Then(assertions func(actual *CapturedIO, givens *InterestingGivens)) *Some {
	assertions(some.capturedIO, some.interestingGivens)
	return some
}

// Call ThenFor in a table test (for loop). Provide a function in which assertions will be made.
// Use the TestingT typed var in place of testing.T.
// The test state is recorded in TestingT type and goGiven fails the test if the error methods (ErrorF etc)
// were called after the function exits.
func (some *Some) ThenFor(assertions func(testingT TestingT, actual *CapturedIO, givens *InterestingGivens)) *Some {
	assertions(some.testingT, some.capturedIO, some.interestingGivens)
	if some.testingT.Failed {
		globalTestingT := some.globalTestingT
		globalTestingT.Helper()
		globalTestingT.Errorf(some.testingT.TestOutput)
	}
	return some
}
