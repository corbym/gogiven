package gogiven

import (
	"io/ioutil"
)

//GivenContext is the interface that exports the internals of TestContext
type GivenContext interface {
	SomeTests() map[string]interface{}
	FileName() string
}

//TestContext contains a safeMap of the TestMetaData for the current test file being processed and
// a copy of the fileName with it's file content.
type TestContext struct {
	someTests   *safeMap
	fileName    string
	fileContent string
}

//NewTestContext creates a new context. This will read the whole contents of filename
// in using ioutil.ReadFile into memory.
func NewTestContext(fileName string) *TestContext {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("file not found:" + err.Error())
	}
	context := new(TestContext)
	context.someTests = newSafeMap()
	context.fileName = fileName
	context.fileContent = string(content[:])
	return context
}

//SomeTests is a map containing the TestMetaData for this TestContext's tests
// that are being executed.
func (c *TestContext) SomeTests() *safeMap {
	return c.someTests
}