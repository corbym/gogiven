package gogiven

import (
	"testing"
)

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
func newInterestingGivens() *InterestingGivens {
	givens := new(InterestingGivens)
	givens.Givens = map[string]interface{}{}
	return givens
}
func newCapturedIO() *CapturedIO {
	capturedIO := new(CapturedIO)
	capturedIO.CapturedIO = map[string]interface{}{}
	return capturedIO
}

func (some *some) When(action ...func(actual *CapturedIO, givens *InterestingGivens)) *some {
	action[0](some.CapturedIO, some.InterestingGivens) // TODO: there could be multiple actions..
	return some
}

var finished = false

func (some *some) Then(assertions func(testingT *TestingT, actual *CapturedIO, givens *InterestingGivens)) *some {
	assertions(some.testingT, some.CapturedIO, some.InterestingGivens)
	return some
}

func newTestMetaData(t *testing.T, testName string) *TestingT {
	testContext := new(TestingT)
	testContext.t = t
	testContext.TestName = testName
	return testContext
}
