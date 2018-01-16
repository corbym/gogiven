package gogiven_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/corbym/gocrest"
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	. "github.com/corbym/gogiven"
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/testdata"
	"strings"
)

type StubHtmlGenerator struct{}

func (*StubHtmlGenerator) Generate(testContext *TestContext) string {
	return "testing"
}

func (*StubHtmlGenerator) FileExtension() string {
	return ".html"
}

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
	var testingT = new(base.TestMetaData)
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
		given := Given(testingT, aFakeGenerator)
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
			Given(tInner).
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

func TestParseGivenWhenThen_Panics(t *testing.T) {
	defer func() {
		rcv := recover()
		AssertThat(t, rcv, is.Not(is.Nil()))
	}()
	ParseGivenWhenThen("foo", "Arfg")
}

const expectedOutput = `Given testing some Data Setup 
			 When something Happens 
			 Then 
		
				 Assert That testing actual " actual " is Equal To " some output "`

func TestParseGivenWhenThen_TextOutput(testing *testing.T) {
	testing.Skip("skipped because whitespace")
	givenWhenThen := ParseGivenWhenThen(".TestMyFirst", "./example_test.go")

	AssertThat(testing, strings.TrimSpace(givenWhenThen), is.EqualTo(expectedOutput))
}

func TestParseGivenWhenThen_RangedTextOutput(testing *testing.T) {
	testing.Skip("skipped because whitespace")
	givenWhenThen := ParseGivenWhenThen(".TestMyFirst_Ranged", "./example_test.go")
	AssertThat(testing, givenWhenThen, is.EqualTo(
		`Given testing some Data Setup 
				 givens " actual " test actual 
		
				 When some Action test 
				 Then 
		
					 Assert That t givens " actual " has Length test expected`))

}

func fileExists(pathToFile string) interface{} {
	fileInfo, err := os.Stat(pathToFile)
	if err != nil {
		return err
	}
	return fileInfo
}

func inTmpDir() *gocrest.Matcher {
	matcher := new(gocrest.Matcher)
	matcher.Matches = func(actual interface{}) bool {
		file, ok := actual.(os.FileInfo)
		if ok {
			matcher.Describe = fmt.Sprintf("%s", file.Name())
			return true
		}
		return false
	}
	return matcher
}

func aFakeGenerator(givens testdata.InterestingGivens) {
	//Generator = new(StubHtmlGenerator) // you too can override Generator and generate any kind of file output.
}

func someDataSetup(givens testdata.InterestingGivens) {
	givens["1"] = "hi"
	givens["2"] = "foo"
	aFakeGenerator(givens)
}

func andMoreDataSetup(givens testdata.InterestingGivens) {
	givens["blarg"] = 12
}

func someAction(capturedIo testdata.CapturedIO, givens testdata.InterestingGivens) {
	capturedIo["foo"] = "foob"
}

func ofFileInTmpDir(fileName string) string {
	return fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, fileName)
}
