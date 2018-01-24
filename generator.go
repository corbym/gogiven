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
var OutputListeners = []generator.OutputListener{new(generator.FileOutputGenerator)}

func transformFileNameToHeader(fileName string) (header string) {
	return strings.Title(strings.Replace(strings.TrimSuffix(filepath.Base(fileName), ".go"), "_", " ", -1))
}

// GenerateTestOutput generates the test output. Call this method from TestMain.
func GenerateTestOutput() {
	for _, key := range globalTestContextMap.Keys() {
		value, _ := globalTestContextMap.Load(key)
		currentTestContext := value.(*TestContext)
		tests := currentTestContext.SomeTests()

		pageData := generator.NewPageData(
			transformFileNameToHeader(currentTestContext.fileName),
			tests.asMapOfSome(),
		)

		output := Generator.Generate(pageData)
		contentType := Generator.ContentType()
		for _, listener := range OutputListeners {
			listener.Notify(currentTestContext.fileName, contentType, output)
		}
	}
}
