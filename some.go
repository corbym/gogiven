package gogiven

type some struct {
	testingT          *TestMetaData
	InterestingGivens *InterestingGivens
	CapturedIO        *CapturedIO
}

func newSome(testContext *TestMetaData,
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
	action[0](some.CapturedIO, some.InterestingGivens) // TODO: there could be multiple actions..
	return some
}

func (some *some) Then(assertions func(actual *CapturedIO, givens *InterestingGivens)) *some {
	assertions(some.CapturedIO, some.InterestingGivens)
	return some
}
