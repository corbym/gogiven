package gogiven

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/testdata"
	"testing"
)

func TestMyFirst(testing *testing.T) {
	Given(testing, theSystemSetup).
		When(somethingHappens).
		Then(func(testing base.TestingT,
			actual testdata.CapturedIO,
			givens testdata.InterestingGivens,
		) { // passed in testing should be used for assertions

			//we do some assertions here, commenting why
			AssertThat(testing, actual["actual"], is.EqualTo("some output"))
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
		{actual: "a", expected: 1},
	}
	for _, test := range someRange {
		Given(testing, theSystemSetup, func(givens testdata.InterestingGivens) {
			givens["actual"] = test.actual
		}).
			When(somethingHappensWithThe(test)).
			Then(func(t base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
				//do assertions
				AssertThat(t, givens["actual"], has.Length(test.expected))
			})
	}
}

func TestMyFirst_Skipped(tst *testing.T) {
	var someRange = []struct {
		actual   string
		expected int
	}{
		{actual: "fff", expected: 0},
		{actual: "a", expected: 1},
	}
	for _, test := range someRange {
		tst.Run(test.actual, func(weAreTesting *testing.T) {
			Given(weAreTesting, theSystemSetup, func(givens testdata.InterestingGivens) {
				givens["actual"] = test.actual
			}).
				SkippingThisOneIf(func(someData ...interface{}) bool {
					return test.actual == "fff"
				}, "some data %s does not work yet", test.actual).
				When(somethingHappensWithThe(test)).
				Then(func(t base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
					AssertThat(t, test.actual, is.EqualTo("a").Reason("we only want to assert if test actual is a not empty"))
				})
		})
	}
}

func TestWithoutGiven(t *testing.T) {
	When(t, somethingHappens).
		Then(func(testing base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
			AssertThat(testing, actual["actual"], is.EqualTo("some output"))
		})
}

func somethingHappensWithThe(data struct {
	actual   string
	expected int
}) base.CapturedIOGivenData {
	return func(capturedIO testdata.CapturedIO, givens testdata.InterestingGivens) {
		capturedIO[data.actual] = data.expected
	}
}

func theSystemSetup(givens testdata.InterestingGivens) {
	givens["foofar"] = "faff"
}
