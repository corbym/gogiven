package gogiven_test

import (
	"fmt"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	. "github.com/corbym/gogiven"
	"path/filepath"
	"testing"
)

var rawContent string
var html string

func TestGeneratorCreatesTestFirstTestHtml(testing *testing.T) {

	Given(testing, someFileContent).
		When(fileIsConvertedToHtml).
		Then(func(testing TestingT, actual *CapturedIO, givens *InterestingGivens) {

		AssertThat(testing, html, is.ValueContaining("<title>Given Test</title>"))
		AssertThat(testing, html, is.ValueContaining("<h1>Given Test</h1>"))
		AssertThat(testing, html, is.ValueContaining("Given testing some Data Setup"))
		AssertThat(testing, html, is.ValueContaining("Given testing some Data Setup and More Data Setup"))
	})
}

func fileIsConvertedToHtml(actual *CapturedIO, givens *InterestingGivens) {
	context := NewTestContext(testFileName())

	someTests := context.SomeTests()
	for _, funcName := range []string{
		"github.com.corbym.gogiven_test.TestGivenWhenGeneratesHtml",
		"github.com.corbym.gogiven_test.TestGivenWhenStacksGivens",
	} {
		storeTests(someTests, funcName)
	}

	html = Generator.Generate(context)
}

func storeTests(someTests *SafeMap, funcName string) {
	globalTestingT := NewTestMetaData(funcName)
	testMetaData := NewTestMetaData(funcName)
	some := NewSome(globalTestingT,
		funcName,
		testMetaData,
		ParseGivenWhenThen(funcName, fileContent(testFileName())),
		func(givens *InterestingGivens) {
			givens.Givens["foofar"] = "farfoo"
		})
	someTests.Store(funcName, some)
}
func someFileContent(givens *InterestingGivens) {
	rawContent = fileContent(testFileName())
	givens.Givens["file content"] = rawContent
}

func testFileName() string {
	return fmt.Sprintf(".%c%s", filepath.Separator, "given_test.go")
}
