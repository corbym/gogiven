package gogiven

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/testdata"
	"testing"
)

func TestMyFirst(testing *testing.T) {
	Given(testing, someDataSetup).
		When(somethingHappens).
		Then(func(testing base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) { // passed in testing should be used for assertions
			//do assertions
			then.AssertThat(testing, actual["actual"], is.EqualTo("some output"))
		})
}

func somethingHappens(actual testdata.CapturedIO, expected testdata.InterestingGivens) {
	actual["actual"] = "some output"
}

func TestMyFirst_Ranged(testing *testing.T) {
	var someRange = []struct {
		actual   string
		expected int
	}{
		{actual: "", expected: 0},
		{actual: "a", expected: 2},
	}
	for _, test := range someRange {
		Given(testing, someDataSetup, func(givens testdata.InterestingGivens) {
			givens["actual"] = test.actual
		}).
			When(someAction(test)).
			Then(func(t base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
				//do assertions
				then.AssertThat(t, givens["actual"], has.Length(test.expected))
			})
	}
}

func someAction(data struct {
	actual   string
	expected int
}) base.CapturedIOGivenData {
	return func(capturedIO testdata.CapturedIO, givens testdata.InterestingGivens) {
		capturedIO[data.actual] = data.expected
	}
}

func someDataSetup(givens testdata.InterestingGivens) {
	givens["foofar"] = "faff"
}
