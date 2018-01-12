package gogiven

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GoGivensOutputGenerator is an interface that can be implemented by anything that can generate file content to be output
// after a test has completed.
type GoGivensOutputGenerator interface {
	//Called from GenerateTestOutput
	Generate(context *TestContext) (html string)
	//File Extension for the output generated, e.g ".html"
	FileExtension() string
}

//TestOutputGenerator is an implmentation of the GoGivensOutputGenerator that generates an html file per
// test.
type TestOutputGenerator struct {
	GoGivensOutputGenerator
}

// FileExtension for the output generated.
func (generator *TestOutputGenerator) FileExtension() string {
	return ".html"
}

//PageData is the struct that populates the template with data from the test output.
type PageData struct {
	Title   string
	SomeMap map[string]*Some
}

// Generate generates the default output for a test. The return string contains the html
// that goes into the output file generated in gogivens.GenerateTestOutput()
func (generator *TestOutputGenerator) Generate(context *TestContext) string {
	goPath := os.Getenv("GOPATH")

	tmpl := template.Must(template.ParseFiles(
		filepath.Join(goPath, "src/github.com/corbym/gogiven/resources/htmltemplate.gtl"),
		filepath.Join(goPath, "src/github.com/corbym/gogiven/resources/capturedio.gtl"),
		filepath.Join(goPath, "src/github.com/corbym/gogiven/resources/interestinggivens.gtl"),
		filepath.Join(goPath, "src/github.com/corbym/gogiven/resources/style.gtl"),
		filepath.Join(goPath, "src/github.com/corbym/gogiven/resources/test-body.gtl"),
		filepath.Join(goPath, "src/github.com/corbym/gogiven/resources/contents.gtl"),
		filepath.Join(goPath, "src/github.com/corbym/gogiven/resources/javascript.gtl"),
	))
	safeMap := context.someTests
	var buffer bytes.Buffer
	pageData := &PageData{
		Title:   TransformFileNameToHeader(context.fileName),
		SomeMap: safeMap.AsMapOfSome(),
	}
	err := tmpl.ExecuteTemplate(&buffer, "base", pageData)
	if err != nil {
		panic(err.Error())
	}
	return buffer.String()
}

// TransformFileNameToHeader takes a test filename e.g. /foo/bar/my_test.go and returns a header e.g. "My Test".
// Strips off the file path and removes the extension.
func TransformFileNameToHeader(fileName string) (header string) {
	return strings.Title(strings.Replace(strings.TrimSuffix(filepath.Base(fileName), ".go"), "_", " ", -1))
}

// Generator is a global variable that holds the GoGivensOutputGenerator.
// You can replace the generator with your own if you match the interface here
// and set Generator = new(myFooGenerator).
// Don't forget to add the call to the generator function in a "func TestMain(testing.M)" method
// in your test package.
// One file per test file will be generated containing output.
var Generator GoGivensOutputGenerator = new(TestOutputGenerator)

// GenerateTestOutput generates the test output. Call this method from TestMain.
func GenerateTestOutput() {
	for _, key := range globalTestContextMap.Keys() {
		value, _ := globalTestContextMap.Load(key)
		currentTestContext := value.(*TestContext)

		output := Generator.Generate(currentTestContext)
		extension := Generator.FileExtension()
		outputFileName := fmt.Sprintf("%s%c%s", outputDirectory(),
			os.PathSeparator,
			strings.Replace(filepath.Base(currentTestContext.fileName), ".go", extension, 1))

		err := ioutil.WriteFile(outputFileName, []byte(output), 0644)
		if err != nil {
			panic("error generating gogiven output:" + err.Error())
		}
		fmt.Printf("\ngenerated test output: file://%s\n", strings.Replace(outputFileName, "\\", "/", -1))
	}
}
func outputDirectory() string {
	outputDir := os.Getenv("GOGIVENS_OUTPUT_DIR")
	if outputDir == "" {
		os.Stdout.WriteString("env var GOGIVENS_OUTPUT_DIR was not found, using TempDir " + os.TempDir())
		outputDir = os.TempDir()
	}
	if _, err := os.Stat(outputDir); err == nil {
		return outputDir
	}
	os.Stderr.WriteString("output dir not found:" + outputDir + ", defaulting to ./")
	return "."
}
