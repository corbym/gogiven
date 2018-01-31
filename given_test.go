package gogiven_test

import (
	"os"
	"testing"

	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	. "github.com/corbym/gogiven"
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/testdata"
)

func TestMain(testmain *testing.M) {
	runOutput := testmain.Run()
	GenerateTestOutput()
	os.Exit(runOutput)
}

func TestGivenWhenSetsInterestingGiven(testing *testing.T) {
	testing.Parallel()
	Given(testing, someDataSetup).
		When(someAction).
		Then(func(t base.TestingT,
			actual testdata.CapturedIO,
			givens testdata.InterestingGivens) {
			//do assertions
			AssertThat(t, actual["foo"], is.EqualTo("foob"))
		})
}

func TestGivenWhenExercisingRanges(testing *testing.T) {
	var testMetaData []*base.TestMetaData
	var testingT = &base.TestMetaData{TestID: "title"}

	var some []*base.Some
	testing.Parallel()
	var someRange = []struct {
		actual   string
		expected int
	}{
		{actual: "", expected: 0},
		{actual: "a", expected: 2},
	}
	for _, test := range someRange {
		given := Given(testingT)
		some = append(some, given)
		given.When(func(actual testdata.CapturedIO, givens testdata.InterestingGivens) {
			actual["value"] = test.actual
			actual["expected"] = test.expected
		}).
			Then(func(t base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
				//do assertions
				testMetaData = append(testMetaData, t.(*base.TestMetaData))
				AssertThat(t, test.actual, has.Length(test.expected))
			})
	}
	AssertThat(testing, some[0].CapturedIO()["value"], is.EqualTo(""))
	AssertThat(testing, some[0].CapturedIO()["expected"], is.EqualTo(0))
	AssertThat(testing, some[1].CapturedIO()["value"], is.EqualTo("a"))
	AssertThat(testing, some[1].CapturedIO()["expected"], is.EqualTo(2))
	AssertThat(testing, testMetaData[0].Failed(), is.EqualTo(false))
	AssertThat(testing, testMetaData[1].Failed(), is.EqualTo(true))
}

func TestInnerTestRangesOverValues(t *testing.T) {
	var someRange = []struct {
		value    string
		expected int
	}{
		{value: "n", expected: 1},
		{value: "aa", expected: 2},
	}
	for _, test := range someRange {
		t.Run(test.value, func(tInner *testing.T) {

			Given(tInner, givenSome_Stuffs("a")).
				When(func(capturedIO testdata.CapturedIO, givens testdata.InterestingGivens) {
					givens["value"] = test.value
					givens["expected"] = test.expected
					capturedIO["actual"] = len(test.value)
				}).Then(func(testingT base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
				AssertThat(testingT, actual["actual"], is.EqualTo(test.expected))
			})
		})
	}
}
func givenSome_Stuffs(myStuff string) base.GivenData {
	return func(givens testdata.InterestingGivens) {

	}
}

func TestGivenWhenStacksGivens(testing *testing.T) {
	testing.Parallel()
	Given(testing, someDataSetup, andMoreDataSetup).
		When(someAction).
		Then(func(testing base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
			//do assertions
			AssertThat(testing, givens, has.AllKeys("1", "2", "blarg"))
			AssertThat(testing, givens, is.ValueContaining("hi", 12, "foo"))
			AssertThat(testing, actual, has.Key("foo"))
		})
}

func TestGivenWhenSkips(testing *testing.T) {
	testing.Parallel()
	t := &base.TestMetaData{TestID: "skiptest"}
	Given(t, someDataSetup, andMoreDataSetup).
		SkippingThisOne("some reason").
		When(someAction).
		Then(func(testing base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
			//do assertions
		})
	AssertThat(testing, t.Skipped(), is.EqualTo(true))
	AssertThat(testing, t.TestOutput(), is.EqualTo("some reason"))
}

func someDataSetup(givens testdata.InterestingGivens) {
	givens["1"] = "hi"
	givens["2"] = "foo"
}

func andMoreDataSetup(givens testdata.InterestingGivens) {
	givens["blarg"] = 12
}

func someAction(capturedIo testdata.CapturedIO, givens testdata.InterestingGivens) {
	capturedIo["foo"] = "foob"
}
