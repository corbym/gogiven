package gogiven_test

import (
	"testing"
	. "github.com/corbym/gogiven"
	. "github.com/corbym/gocrest/then"
	"github.com/corbym/gocrest/is"
	"path/filepath"
	"fmt"
)

var rawContent string
var html string

func TestGeneratorCreatesTestFirstTestHtml(testing *testing.T) {
	Given(testing, someFileContent).
		When(fileIsConvertedToHtml).
		Then(func(t *TestingT, actual *CapturedIO, givens *InterestingGivens) {

		AssertThat(t, html, is.ValueContaining(
			"<html><title>Given Test</title>"+
				"<body>"))
		AssertThat(t, html, is.ValueContaining("<h1>Given Test</h1>"))
		AssertThat(t, html, is.ValueContaining(
			"<div id=\"TestGivenWhenGeneratesHtml\">GIVEN testing some data setup<br/>"+
				"WHEN some action<br/>"+
				"THEN<br/>"+
				"//do assertions<br/>"+
				"Assert that testing \"foo\" is equal to \"foob\"<br/>"+
				"Assert that testing  file exists \"gogiven_test.TestGivenWhenGeneratesHtml.html\" in tmp dir<br/>"+
				"Assert that testing file content of file in tmp dir \"gogiven_test.TestGivenWhenGeneratesHtml.html\" is equalto \"testing\"<br/>"+
				"</div>"))
		AssertThat(t, html, is.ValueContaining(
			"</body>"+
				"</html>"))
	})
}

func fileIsConvertedToHtml(actual *CapturedIO, givens *InterestingGivens) {
	html = Generator.Generate(
		"C:\\Users\\Matt\\go\\src\\github.com\\corbym\\gogiven\\given_test.go",
		rawContent,
	)
	actual.CapturedIO["html"] = html
}

func someFileContent(givens *InterestingGivens) {
	rawContent = fileContent(fmt.Sprintf(".%s%s", filepath.Separator, "given_test.go"))
	givens.Givens["file content"] = rawContent
}
