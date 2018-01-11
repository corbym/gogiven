package gogiven_test

import (
	"fmt"
	"github.com/corbym/gocrest"
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	. "github.com/corbym/gogiven"
	"io/ioutil"
	"os"
	"testing"
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

func TestGivenWhenGeneratesHtml(testing *testing.T) {
	Given(testing, someDataSetup).
		When(someAction).
		Then(func(t TestingT, actual *CapturedIO, givens *InterestingGivens) {
			//do assertions
			AssertThat(t, actual.CapturedIO["foo"], is.EqualTo("foob"))
		})

	AssertThat(testing, fileExists("given_test.html"), inTmpDir())
	AssertThat(testing, fileContent(ofFileInTmpDir("given_test.html")), is.EqualTo("testing"))
}

func TestGivenWhenExercisingRanges(testing *testing.T) {
	var testMetaData []*TestMetaData
	var testingT = new(TestMetaData)
	var some []*Some

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
		given.When(func(actual *CapturedIO, givens *InterestingGivens) {
			actual.CapturedIO["actual"] = test.actual
			actual.CapturedIO["expected"] = test.expected
		}).
			Then(func(t TestingT, actual *CapturedIO, givens *InterestingGivens) {
				//do assertions
				testMetaData = append(testMetaData, t.(*TestMetaData))
				AssertThat(t, test.actual, has.Length(test.expected))
			})
	}
	AssertThat(testing, some[0].CapturedIO()["actual"], is.EqualTo(""))
	AssertThat(testing, some[0].CapturedIO()["expected"], is.EqualTo(0))
	AssertThat(testing, some[1].CapturedIO()["actual"], is.EqualTo("a"))
	AssertThat(testing, some[1].CapturedIO()["expected"], is.EqualTo(2))
	AssertThat(testing, testMetaData[0].Failed(), is.EqualTo(false))
	AssertThat(testing, testMetaData[1].Failed(), is.EqualTo(true))
}

func TestGivenWhenStacksGivens(testing *testing.T) {
	Given(testing, someDataSetup, andMoreDataSetup).
		When(someAction).
		Then(func(testing TestingT, actual *CapturedIO, givens *InterestingGivens) {
			//do assertions
			AssertThat(testing, givens.Givens, has.AllKeys("1", "2", "blarg"))
			AssertThat(testing, givens.Givens, is.ValueContaining("hi", 12, "foo"))
			AssertThat(testing, actual.CapturedIO, has.Key("foo"))
		})
}

func fileExists(fileName string) interface{} {
	dir := ofFileInTmpDir(fileName)
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}
	return fileInfo
}

func fileContent(fileName string) string {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ""
	}
	return string(dat[:])
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

func aFakeGenerator(givens *InterestingGivens) {
	Generator = new(StubHtmlGenerator) // you too can override Generator and generate any kind of file output.
}

func someDataSetup(givens *InterestingGivens) {
	givens.Givens["1"] = "hi"
	givens.Givens["2"] = "foo"
	aFakeGenerator(givens)
}

func andMoreDataSetup(givens *InterestingGivens) {
	givens.Givens["blarg"] = 12
}

func someAction(capturedIo *CapturedIO, givens *InterestingGivens) {
	capturedIo.CapturedIO["foo"] = "foob"
}

func ofFileInTmpDir(fileName string) string {
	return fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, fileName)
}
