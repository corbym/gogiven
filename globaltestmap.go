package gogiven

import (
	"io/ioutil"
)

type GivenContext interface {
	SomeTests() map[string]interface{}
	FileName() string
	FileContent() string
}

type TestContext struct {
	someTests   *SafeMap
	fileName    string
	fileContent []byte
}

func NewTestContext(fileName string) *TestContext {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("file not found:" + err.Error())
	}
	context := new(TestContext)
	context.someTests = newSafeMap()
	context.fileName = fileName
	context.fileContent = content
	return context
}

func (c *TestContext) FileName() string {
	return c.fileName
}

func (c *TestContext) SomeTests() *SafeMap {
	return c.someTests
}
func (c *TestContext) FileContent() []byte {
	return c.fileContent
}
