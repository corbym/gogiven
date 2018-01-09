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
		Then(func(actual *CapturedIO, givens *InterestingGivens) {

		AssertThat(testing, html, is.ValueContaining(
			"<html><title>Given Test</title>"+
				"<body>"))
		AssertThat(testing, html, is.ValueContaining("<h1>Given Test</h1>"))
		AssertThat(testing, html, is.ValueContaining(
			"<div id=\"TestGivenWhenGeneratesHtml\">GIVEN testing some data setup<br/>"+
				"WHEN some action<br/>"+
				"THEN<br/>"+
				"//do assertions<br/>"+
				"Assert that testing \"foo\" is equal to \"foob\"<br/>"+
				"Assert that testing  file exists \"gogiven_test.TestGivenWhenGeneratesHtml.html\" in tmp dir<br/>"+
				"Assert that testing file content of file in tmp dir \"gogiven_test.TestGivenWhenGeneratesHtml.html\" is equalto \"testing\"<br/>"+
				"</div>"))
		AssertThat(testing, html, is.ValueContaining(
			"</body>"+
				"</html>"))
	})
}

func fileIsConvertedToHtml(actual *CapturedIO, givens *InterestingGivens) {
	context := NewTestContext(testFileName())

	someTests := context.SomeTests()
	testName := "TestGivenWhenGeneratesHtml"
	globalTestingT := NewTestMetaData(testName)
	testMetaData := NewTestMetaData(testName)
	some := NewSome(globalTestingT,
		testMetaData,
		ParseGivenWhenThen(testName, fileContent(testFileName())),
		func(givens *InterestingGivens) {
			givens.Givens["foofar"] = "farfoo"
		})
	someTests.Store(testName, some)

	html = Generator.Generate(context)
}
func someFileContent(givens *InterestingGivens) {
	rawContent = fileContent(testFileName())
	givens.Givens["file content"] = rawContent
}
func testFileName() string {
	return fmt.Sprintf(".%c%s", filepath.Separator, "given_test.go")
}
