package gogiven_test

import (
	"testing"
	"os"
	"fmt"
	. "github.com/corbym/gocrest/then"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest"
	. "github.com/corbym/gogiven/given"
)

func TestGivenWhenGeneratesHtml(testing *testing.T) {
	Given(testing, someDataSetup).
		When(someAction).
		Then(func(actual *CapturedIO, givens *InterestingGivens) {
		//do assertions
		AssertThat(testing, actual.CapturedIO["foo"], is.EqualTo("foob"))
	})

	AssertThat(testing, fileExists("gogiven_test.TestGivenWhenGeneratesHtml.html"), inTmpDir())
}

func TestGivenWhenStacksGivens(testing *testing.T) {
	Given(testing, someDataSetup, andMoreDataSetup).
		When(someAction).
		Then(func(actual *CapturedIO, givens *InterestingGivens) {
		//do assertions
		AssertThat(testing, givens.Givens, has.AllKeys("1", "2", "blarg"))
		AssertThat(testing, givens.Givens, is.ValueContaining("hi", 12, "foo"))
		AssertThat(testing, actual, is.EqualTo("foo"))
	})
}

func fileExists(fileName string) (error) {
	_, err := os.Stat(fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, fileName))
	return err
}

func inTmpDir() *gocrest.Matcher {
	matcher := new(gocrest.Matcher)
	matcher.Matches = func(actual interface{}) bool {
		matcher.Describe = fmt.Sprintf("%s", actual.(error).Error())
		fileError := actual.(error)
		if os.IsNotExist(fileError) {
			return false
		}
		return false
	}
	return matcher
}

func someDataSetup(givens *InterestingGivens) *InterestingGivens {
	givens.Givens["1"] = "hi"
	givens.Givens["2"] = "foo"
	return givens
}

func andMoreDataSetup(givens *InterestingGivens) *InterestingGivens {
	givens.Givens["blarg"] = 12
	return givens
}

func someAction(capturedIo *CapturedIO, givens *InterestingGivens) *CapturedIO {
	capturedIo.CapturedIO["foo"] = "foob"
	return capturedIo
}
