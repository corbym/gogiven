package gogiven

import "github.com/corbym/gocrest"

type some struct {
	globalTestingT    gocrest.TestingT
	testingT          *TestMetaData
	InterestingGivens *InterestingGivens
	CapturedIO        *CapturedIO
}

func newSome(
	globalTestingT gocrest.TestingT,
	testContext *TestMetaData,
	givenFunc ...func(givens *InterestingGivens)) *some {

	some := new(some)
	some.testingT = testContext
	some.globalTestingT = globalTestingT
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

func (some *some) ThenFor(assertions func(testingT *TestMetaData, actual *CapturedIO, givens *InterestingGivens)) *some {
	assertions(some.testingT, some.CapturedIO, some.InterestingGivens)
	if some.testingT.Failed {
		globalTestingT := some.globalTestingT
		globalTestingT.Helper()
		globalTestingT.Errorf(some.testingT.TestOutput)
	}
	return some
}
