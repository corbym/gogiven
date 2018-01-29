package gogiven

import (
	"github.com/corbym/gogiven/generator"
	"github.com/corbym/htmlspec"
	"path/filepath"
	"strings"
)

// Generator is a global variable that holds the GoGivensOutputGenerator.
// You can replace the generator with your own if you match the interface here
// and set Generator = new(myFooGenerator) in a method (usually TestMain or init).
// Don't forget to add the call to the generator function in a "func TestMain(testing.M)" method
// in your test package.
// One file per test file will be generated containing output.
var Generator generator.GoGivensOutputGenerator = htmlspec.NewTestOutputGenerator()

//OutputListeners holds a list of listners which can process the output generated by the GoGivensOutputGenerator.
// By default, the generator.FileOutputGenerator is added. but you can append or replace with your own.
var OutputListeners = []generator.OutputListener{new(generator.FileOutputGenerator)}

func transformFileNameToHeader(fileWithPath string) (header string) {
	baseFileName := filepath.Base(fileWithPath)
	withTrimmedSuffix := strings.TrimSuffix(baseFileName, ".go")
	return strings.Title(strings.Replace(withTrimmedSuffix, "_", " ", -1))
}

// GenerateTestOutput generates the test output. Call this method from TestMain.
// The global var Generator is used to generate the content, and the OutputListeners are iterated and Notified
// of the content that has been generated.
func GenerateTestOutput() {
	for _, key := range globalTestContextMap.Keys() {
		value, _ := globalTestContextMap.Load(key)
		currentTestContext := value.(*TestContext)
		tests := currentTestContext.SomeTests()

		pageData := generator.NewPageData(
			transformFileNameToHeader(currentTestContext.FileName()),
			tests.asMapOfSome(),
		)

		output := Generator.Generate(pageData)
		contentType := Generator.ContentType()
		for _, listener := range OutputListeners {
			listener.Notify(currentTestContext.FileName(), contentType, output)
		}
	}
}
