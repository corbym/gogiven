package gogiven

type TestingT interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
	Name() string
}

type Some struct {
	globalTestingT    TestingT
	testingT          *TestMetaData
	interestingGivens *InterestingGivens
	capturedIO        *CapturedIO
}

func newSome(
	globalTestingT TestingT,
	testContext *TestMetaData,
	givenFunc ...func(givens *InterestingGivens)) *Some {

	some := new(Some)
	some.testingT = testContext
	some.globalTestingT = globalTestingT
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
func (some *Some) InterestingGivens() map[string]interface{} {
	return some.interestingGivens.Givens
}

func (some *Some) CapturedIO() map[string]interface{} {
	return some.capturedIO.CapturedIO
}

func (some *Some) When(action ...func(actual *CapturedIO, givens *InterestingGivens)) *Some {
	action[0](some.capturedIO, some.interestingGivens) // TODO: there could be multiple actions..
	return some
}

func (some *Some) Then(assertions func(actual *CapturedIO, givens *InterestingGivens)) *Some {
	assertions(some.capturedIO, some.interestingGivens)
	return some
}

func (some *Some) ThenFor(assertions func(testingT *TestMetaData, actual *CapturedIO, givens *InterestingGivens)) *Some {
	assertions(some.testingT, some.capturedIO, some.interestingGivens)
	if some.testingT.Failed {
		globalTestingT := some.globalTestingT
		globalTestingT.Helper()
		globalTestingT.Errorf(some.testingT.TestOutput)
	}
	return some
}