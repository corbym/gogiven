package gogiven_test

import (
	"testing"
	. "github.com/corbym/gogiven"
	"github.com/corbym/gocrest/is"
	"fmt"
	"path/filepath"
	. "github.com/corbym/gocrest/then"
)

var rawContent string
var html string

func TestGeneratorCreatesTestFirstTestHtml(testing *testing.T) {
	testing.Skipf("Need to implement this")
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
	html = Generator.Generate(new(TestContext))
	actual.CapturedIO["html"] = html
}

func someFileContent(givens *InterestingGivens) {
	rawContent = fileContent(fmt.Sprintf(".%s%s", filepath.Separator, "given_test.go"))
	givens.Givens["file content"] = rawContent
}
