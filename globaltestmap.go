package gogiven

import (
	"io/ioutil"
)

type TestContext struct {
	someTests   *SafeMap
	fileName    string
	fileContent []byte
}

func NewGlobalTestContext(fileName string) *TestContext {
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
