package gogiven_test

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven"
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/testdata"
	"testing"
)

func TestMyFirst(t *testing.T) {
	gogiven.Given(t, theSystemSetup).
		When(somethingHappens).
		Then(func(t base.TestingT, captured testdata.CapturedIO, givens testdata.InterestingGivens) {
			// we do some assertions here, commenting why
			AssertThat(t, captured["actual"], is.EqualTo("some output"))
		})
}

func somethingHappens(actual testdata.CapturedIO, expected testdata.InterestingGivens) {
	actual["actual"] = "some output"
}

// This test tests over a range of values. Lorum ipsum dolor, lorum ipsum dolor lorum ipsum dolor. Lorum ipsum dolor.
// Comments on new lines will be split into paragraph sections.
//
// Empty lines will be removed. Do not remove this comment. Thanks.
func TestMyFirst_Ranged(t *testing.T) {
	var someRange = []struct {
		actual   string
		expected int
	}{
		{actual: "x", expected: 2},
		{actual: "aaaa", expected: 4},
	}
	for _, test := range someRange {
		t.Run(test.actual, func(tt *testing.T) {
			weAreTesting := base.NewTestMetaData(t.Name()) // this test is fake, as we want to demo failing
			gogiven.Given(weAreTesting, theSystemSetup, withTestData(test)).
				When(somethingHappensWithThe(test)).
				Then(func(t base.TestingT, captured testdata.CapturedIO, stored testdata.InterestingGivens) {
					// do assertions
					AssertThat(t, stored["actual"], has.Length(test.expected))
				})
		})
	}
}
func withTestData(test someData) func(givens testdata.InterestingGivens) {
	return func(givens testdata.InterestingGivens) {
		givens["actual"] = test.actual
	}
}

func TestMyFirst_Skipped(t *testing.T) {
	var someRange = []struct {
		actual   string
		expected int
	}{
		{actual: "fff", expected: 0},
		{actual: "a", expected: 1},
	}
	for _, test := range someRange {
		t.Run(test.actual, func(t *testing.T) {
			gogiven.Given(t, theSystemSetup, thatIsABitDodgyTo(test)).
				SkippingThisOneIf(theValueIsFff(test), "some data %s does not work yet", test.actual).
				When(somethingHappensWithThe(test)).
				Then(func(t base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
					AssertThat(t, test.actual, is.EqualTo("a").Reason("we only want to assert if test actual is a not empty"))
				})
		})
	}
}
func thatIsABitDodgyTo(test someData) func(givens testdata.InterestingGivens) {
	return func(givens testdata.InterestingGivens) {
		givens["actual"] = test.actual
	}
}
func theValueIsFff(someData someData) func(someData ...interface{}) bool {
	return func(data ...interface{}) bool {
		return someData.actual == "fff"
	}
}

//func theValueIsFff(test ...interface{}) bool {
//	return test[0].(*someData).actual == "fff"
//}

func TestMyFirst_NonDefaultParamName(myT *testing.T) {
	gogiven.Given(myT, theSystemSetup).
		When(somethingHappens).
		Then(func(thenT base.TestingT, captured testdata.CapturedIO, givens testdata.InterestingGivens) {
			// we do some assertions here, commenting why
			AssertThat(thenT, captured["actual"], is.EqualTo("some output"))
		})
}

func TestWithoutGiven(t *testing.T) {
	gogiven.When(t, somethingHappens).
		Then(func(t base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
			AssertThat(t, actual["actual"], is.EqualTo("some output"))
		})
}

type someData struct {
	actual   string
	expected int
}

func somethingHappensWithThe(data someData) base.CapturedIOGivenData {
	return func(capturedIO testdata.CapturedIO, givens testdata.InterestingGivens) {
		capturedIO[data.actual] = data.expected
	}
}

func theSystemSetup(givens testdata.InterestingGivens) {
	givens["foofar"] = "faff"
}
