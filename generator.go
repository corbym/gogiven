package gogiven

import (
	"fmt"
	"github.com/corbym/gogiven/generator"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/corbym/htmlspec"
)

// Generator is a global variable that holds the GoGivensOutputGenerator.
// You can replace the generator with your own if you match the interface here
// and set Generator = new(myFooGenerator) in a method (usually TestMain or init).
// Don't forget to add the call to the generator function in a "func TestMain(testing.M)" method
// in your test package.
// One file per test file will be generated containing output.
var Generator generator.GoGivensOutputGenerator = htmlspec.NewTestOutputGenerator()

func transformFileNameToHeader(fileName string) (header string) {
	return strings.Title(strings.Replace(strings.TrimSuffix(filepath.Base(fileName), ".go"), "_", " ", -1))
}

// GenerateTestOutput generates the test output. Call this method from TestMain.
func GenerateTestOutput() {
	for _, key := range globalTestContextMap.Keys() {
		value, _ := globalTestContextMap.Load(key)
		currentTestContext := value.(*TestContext)
		tests := currentTestContext.SomeTests()
		pageData := &generator.PageData{
			Title:   transformFileNameToHeader(currentTestContext.FileName()),
			SomeMap: tests.AsMapOfSome(),
		}
		output := Generator.Generate(pageData)
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
