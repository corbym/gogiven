package gogiven

import (
	"testing"
)

type TestMetaData struct {
	TestId string
	t      *testing.T
}

func newTestMetaData(t *testing.T, testName string) *TestMetaData {
	testContext := new(TestMetaData)
	testContext.t = t
	testContext.TestId = testName
	return testContext
}
