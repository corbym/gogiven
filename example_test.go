package gogiven

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/then"
	"testing"
)

type ActualExpected struct {
	actual   string
	expected int
}

func TestMyFirst(testing *testing.T) {
	var someRange = []ActualExpected{
		{actual: "", expected: 0},
		{actual: "a", expected: 1},
	}
	for _, test := range someRange {
		Given(testing, someDataSetup, func(givens *InterestingGivens) {
			givens.Givens["actual"] = test.actual
		}).
			When(someAction(test)).
			Then(func(t TestingT, actual *CapturedIO, givens *InterestingGivens) {
				//do assertions
				then.AssertThat(t, test.actual, has.Length(test.expected))
			})
	}
}
func someAction(data ActualExpected) CapturedIOGivenData {
	return func(capturedIO *CapturedIO, givens *InterestingGivens) {
		capturedIO.CapturedIO[data.actual] = data.expected
	}
}

func someDataSetup(givens *InterestingGivens) {
	givens.Givens["foofar"] = "faff"
}
