package gogiven

type some struct {
	testingT          *TestingT
	InterestingGivens *InterestingGivens
	CapturedIO        *CapturedIO
}

func newSome(testContext *TestingT,
	givenFunc ...func(givens *InterestingGivens)) *some {

	some := new(some)
	some.testingT = testContext
	some.CapturedIO = newCapturedIO()
	givens := newInterestingGivens()

	if len(givenFunc) > 0 {
		for _, someGivenFunc := range givenFunc {
			someGivenFunc(givens)
		}
	}
	some.InterestingGivens = givens
	return some
}

func (some *some) When(action ...func(actual *CapturedIO, givens *InterestingGivens)) *some {
	testingT := some.testingT
	testingT.Helper()
	action[0](some.CapturedIO, some.InterestingGivens) // TODO: there could be multiple actions..
	return some
}

func (some *some) Then(assertions func(testingT *TestingT, actual *CapturedIO, givens *InterestingGivens)) *some {
	testingT := some.testingT
	testingT.Helper()
	assertions(some.testingT, some.CapturedIO, some.InterestingGivens)
	return some
}
func (some *some) SkippingThisOne() *some {
	testingT := some.testingT
	testingT.Helper()
	testingT.Skipped()
	return some
}
func (some *some) InParallel() *some {
	testingT := some.testingT
	testingT.Helper()
	testingT.Parallel()
	return some
}