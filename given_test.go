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
	var testingT = &base.TestMetaData{TestId: "title"}

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
			return
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
	t := &base.TestMetaData{TestId: "skiptest"}
	Given(t, someDataSetup, andMoreDataSetup).
		SkippingThisOne("some reason").
		When(someAction).
		Then(func(testing base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) {
			//do assertions
		})
	AssertThat(testing, t.Skipped(), is.EqualTo(true))
	AssertThat(testing, t.TestOutput(), is.EqualTo("some reason"))
}

func TestParseGivenWhenThen_Panics(t *testing.T) {
	defer func() {
		rcv := recover()
		AssertThat(t, rcv, is.Not(is.Nil()))
	}()
	ParseGivenWhenThen("foo", "Arfg")
}

func TestParseGivenWhenThen_TextOutputContent(testing *testing.T) {
	givenWhenThen := ParseGivenWhenThen(".TestMyFirst", "./example_test.go")

	AssertThat(testing, givenWhenThen, has.Length(5))

	AssertThat(testing, givenWhenThen[0], is.EqualTo("Given testing the system setup"))
	AssertThat(testing, givenWhenThen[1], is.EqualTo("When something happens"))
	AssertThat(testing, givenWhenThen[2], is.EqualTo("Then noting that passed in testing should be used for assertions"))
	AssertThat(testing, givenWhenThen[3], is.EqualTo("Noting that we do some assertions here commenting why"))
	AssertThat(testing, givenWhenThen[4], is.EqualTo("Assert that testing actual \"actual\" is equal to \"some output\""))
}

func TestParseGivenWhenThen_WithoutGiven(testing *testing.T) {
	givenWhenThen := ParseGivenWhenThen(".TestWithoutGiven", "./example_test.go")

	AssertThat(testing, givenWhenThen, has.Length(3))

	AssertThat(testing, givenWhenThen[0], is.EqualTo("When t something happens"))
	AssertThat(testing, givenWhenThen[1], is.EqualTo("Then"))
	AssertThat(testing, givenWhenThen[2], is.EqualTo("Assert that testing actual \"actual\" is equal to \"some output\""))
}

func TestParseGivenWhenThen_PanicsWithoutGivenOrWhen(testing *testing.T) {
	defer func() {
		recover := recover()
		AssertThat(testing, recover, is.Not(is.Nil()))
	}()
	ParseGivenWhenThen(".TestParseGivenWhenThen_PanicsWithoutGivenOrWhen", "./given_test.go")
}

func TestParseGivenWhenThen_FuncWithReturnType(testing *testing.T) {
	givenWhenThen := ParseGivenWhenThen(".TestMyFirst_Skipped", "./example_test.go")
	AssertThat(testing, givenWhenThen, has.Length(8))

	AssertThat(testing, givenWhenThen[0], is.EqualTo("Given we are testing the system setup"))
	AssertThat(testing, givenWhenThen[1], is.EqualTo("Givens \"actual\" = test actual"))
	AssertThat(testing, givenWhenThen[2], is.EqualTo("Skipping this one if bool"))
	AssertThat(testing, givenWhenThen[3], is.EqualTo("Test actual == \"fff\""))
	// TODO: fix this
	//AssertThat(testing, givenWhenThen[4], is.EqualTo("some data % s does not work yet \"test actual\""))
}

func TestParseGivenWhenThen_RangedTextOutput(testing *testing.T) {
	givenWhenThen := ParseGivenWhenThen(".TestMyFirst_Ranged", "./example_test.go")
	AssertThat(testing, givenWhenThen, has.Length(6))

	AssertThat(testing, givenWhenThen[0], is.EqualTo("Given testing the system setup"))
	AssertThat(testing, givenWhenThen[1], is.EqualTo("Givens \"actual\" = test actual"))
	AssertThat(testing, givenWhenThen[2], is.EqualTo("When something happens with the test"))
	AssertThat(testing, givenWhenThen[3], is.EqualTo("Then"))
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
