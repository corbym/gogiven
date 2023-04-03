package htmlspec

import (
	"bytes"
	"fmt"
	"github.com/corbym/gogiven/generator"
	"html/template"
	"io"
)

const baseTemplateName = "base"
const indexTemplateName = "index"

// HTMLOutputGenerator is an implementation of the GoGivensOutputGenerator that generates an html file per
// test. It is thread safe between goroutines.
type HTMLOutputGenerator struct {
	generator.GoGivensOutputGenerator
	htmlTemplate  *template.Template
	indexTemplate *template.Template
}

var lastError error

// NewHTMLOutputGenerator NewTestOutputGenerator creates a htmlTemplate that is used to generate the html output.
func NewHTMLOutputGenerator() *HTMLOutputGenerator {
	outputGenerator := new(HTMLOutputGenerator)
	htmlTemplate := parseHTMLTemplate()
	indexTemplate := parseIndexTemplate()

	if lastError != nil {
		panic(lastError.Error())
	}
	outputGenerator.htmlTemplate = htmlTemplate
	outputGenerator.indexTemplate = indexTemplate
	return outputGenerator
}

func parseHTMLTemplate() *template.Template {
	htmlTemplate := template.New(baseTemplateName)
	htmlTemplate.Parse(safeStringConverter(Asset("resources/htmltemplate.gtl")))
	htmlTemplate.Parse(safeStringConverter(Asset("resources/interestinggivens.gtl")))
	htmlTemplate.Parse(safeStringConverter(Asset("resources/capturedio.gtl")))
	htmlTemplate.Parse(safeStringConverter(Asset("resources/style.gtl")))
	htmlTemplate.Parse(safeStringConverter(Asset("resources/test-body.gtl")))
	htmlTemplate.Parse(safeStringConverter(Asset("resources/contents.gtl")))
	htmlTemplate.Parse(safeStringConverter(Asset("resources/javascript.gtl")))
	fmt.Printf("loaded html templates%s\n", htmlTemplate.DefinedTemplates())
	return htmlTemplate
}

func parseIndexTemplate() *template.Template {
	indexTemplate := template.New(indexTemplateName)
	indexTemplate.Funcs(template.FuncMap{"OverallTestResult": OverallTestResult})
	indexTemplate.Parse(safeStringConverter(Asset("resources/javascript.gtl")))
	indexTemplate.Parse(safeStringConverter(Asset("resources/index-content.gtl")))
	indexTemplate.Parse(safeStringConverter(Asset("resources/style.gtl")))
	indexTemplate.Parse(safeStringConverter(Asset("resources/index.gtl")))
	fmt.Printf("loaded index templates%s\n", indexTemplate.DefinedTemplates())
	return indexTemplate
}

func safeStringConverter(asset []byte, err error) string {
	if err != nil {
		lastError = err
	}
	return string(asset[:])
}

// ContentType for the output generated.
func (outputGenerator *HTMLOutputGenerator) ContentType() string {
	return "text/html"
}

// Generate generates html output for a test. The return string contains the html
// that goes into the output generated in gogivens.GenerateTestOutput().
// The function panics if the htmlTemplate cannot be generated.
func (outputGenerator *HTMLOutputGenerator) Generate(pageData generator.PageData) io.Reader {
	var buffer = new(bytes.Buffer)
	outputGenerator.htmlTemplate.ExecuteTemplate(buffer, baseTemplateName, pageData)
	return buffer
}

// GenerateIndex generates an index of all the data from the tests in HTML format.
func (outputGenerator *HTMLOutputGenerator) GenerateIndex(indexData []generator.IndexData) io.Reader {
	var buffer = new(bytes.Buffer)
	err := outputGenerator.indexTemplate.ExecuteTemplate(buffer, indexTemplateName, indexData)
	if err != nil {
		panic(err.Error())
	}
	return buffer
}

func OverallTestResult(testData []generator.TestData) string {
	for _, data := range testData {
		if data.TestResult.Failed {
			return "test-failed"
		}
	}
	return "test-passed"
}
