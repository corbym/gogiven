package gogiven

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Interface that anything that can generate file content to be output
// after a test has completed.
type GoGivensOutputGenerator interface {
	Generate(context *TestContext) (html string)
	FileExtension() string
}

//Implementation of the GoGivensOutputGenerator that generates an html file per
// test.
type TestOutputGenerator struct {
	GoGivensOutputGenerator
}

func (generator *TestOutputGenerator) FileExtension() string {
	return ".html"
}

func (generator *TestOutputGenerator) Generate(context *TestContext) string {
	html := testTemplate(context.fileName, string(context.fileContent[:]))
	safeMap := context.someTests
	for _, key := range safeMap.Keys() {
		if some, ok := safeMap.Load(key); ok {
			tests := some.(*Some)
			html += tests.globalTestingT.Name()
		}
	}
	html += "</body></html>"
	return html
}
func testTemplate(fileName string, testFileContent string) string {
	return testHeader(TransformFileNameToHeader(fileName))
}

func testHeader(title string) string {
	return fmt.Sprintf(
		"<html><title>%s</title>"+
			"<body><h1>%s</h1>", title, title)
}

// Takes a test filename e.g. /foo/bar/my_test.go and returns a header e.g. "My Test".
// Strips off the file path and removes the extension.
func TransformFileNameToHeader(fileName string) (header string) {
	return strings.Title(strings.Replace(strings.TrimSuffix(filepath.Base(fileName), ".go"), "_", " ", -1))
}

// Global variable that holds the GoGivensOutputGenerator.
// You can replace the generator with your own if you match the interface here
// and set GoGivensOutputGenerator = myFooGenerator
var Generator GoGivensOutputGenerator = new(TestOutputGenerator)

func GenerateTestOutput() {
	for _, key := range globalTestContextMap.Keys() {
		value, _ := globalTestContextMap.Load(key)
		currentTestContext := value.(*TestContext)

		output := Generator.Generate(currentTestContext)
		extension := Generator.FileExtension()
		outputFileName := fmt.Sprintf("%s%c%s", os.TempDir(),
			os.PathSeparator,
			strings.Replace(filepath.Base(currentTestContext.fileName), ".go", extension, 1))

		err := ioutil.WriteFile(outputFileName, []byte(output), 0644)
		if err != nil {
			panic("error generating gogiven output:" + err.Error())
		}
		fmt.Printf("\ngenerated test output: file:///%s\n", strings.Replace(outputFileName, "\\", "/", -1))
	}
}
