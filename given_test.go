package gogiven_test

import (
	"fmt"
	"github.com/corbym/gocrest"
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	. "github.com/corbym/gogiven"
	"os"
	"testing"
	"io/ioutil"
)

type StubHtmlGenerator struct{}

func (*StubHtmlGenerator) Generate(fileNameWithPath string, testFileContent string) (html string) {
	return "testing"
}

func init() {
	Generator = new(StubHtmlGenerator)
}

func TestGivenWhenGeneratesHtml(testing *testing.T) {
	var context *TestingT

	Given(testing, someDataSetup).
		When(someAction).
		Then(func(testingT *TestingT, actual *CapturedIO, givens *InterestingGivens) {
		//do assertions
		context = testingT
		AssertThat(testingT, actual.CapturedIO["foo"], is.EqualTo("foob"))
	})

	AssertThat(testing, fileExists("given_test.html"), inTmpDir())
	AssertThat(testing, fileContent(ofFileInTmpDir("given_test.html")), is.EqualTo("testing"))
	AssertThat(testing, context.TestName, is.EqualTo("github.com/corbym/gogiven_test.TestGivenWhenGeneratesHtml"))
	AssertThat(testing, context.HasFailed(), is.EqualTo(false))
}

func TestGivenWhenExercisingRanges(testing *testing.T) {
	var someRange = []struct {
		a int
		b int
	}{
		{1, 2},
		{3, 4},
	}
	for _, test := range someRange {
		Given(testing).
			When(func(actual *CapturedIO, givens *InterestingGivens) {
			actual.CapturedIO[fmt.Sprintf("%d", test.a)] = "fooa"
			actual.CapturedIO[fmt.Sprintf("%d", test.b)] = "foob"
		}).
			Then(func(testContext *TestingT, actual *CapturedIO, givens *InterestingGivens) {
			//do assertions
			AssertThat(testContext, actual.CapturedIO, is.ValueContaining("foob"))
			AssertThat(testContext, actual.CapturedIO, is.ValueContaining("fooa"))
		})
	}
	AssertThat(testing, fileExists("given_test.html"), inTmpDir())
	//AssertThat(testing, context.TestName, is.EqualTo("github.com/corbym/gogiven_test.TestGivenWhenGeneratesHtml"))
	//AssertThat(testing, context.TestName, is.EqualTo("github.com/corbym/gogiven_test.TestGivenWhenGeneratesHtml_1"))
}

func TestGivenWhenStacksGivens(testing *testing.T) {
	Given(testing, someDataSetup, andMoreDataSetup).
		When(someAction).
		Then(func(testContext *TestingT, actual *CapturedIO, givens *InterestingGivens) {
		//do assertions
		AssertThat(testContext, givens.Givens, has.AllKeys("1", "2", "blarg"))
		AssertThat(testContext, givens.Givens, is.ValueContaining("hi", 12, "foo"))
		AssertThat(testContext, actual.CapturedIO, has.Key("foo"))
	})
}

func fileExists(fileName string) (interface{}) {
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

func someDataSetup(givens *InterestingGivens) {
	givens.Givens["1"] = "hi"
	givens.Givens["2"] = "foo"
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
