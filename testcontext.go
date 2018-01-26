package gogiven

//TestContext contains a safeMap of the TestMetaData for the current test file being processed.//////
type TestContext struct {
	someTests *safeMap
	fileName  string
}

//NewTestContext creates a new context.
func NewTestContext(fileName string) *TestContext {
	context := new(TestContext)
	context.someTests = newSafeMap()
	context.fileName = fileName
	return context
}

//SomeTests is a map containing the TestMetaData for this TestContext's tests
// that are being executed.
func (c *TestContext) SomeTests() *safeMap {
	return c.someTests
}

//FileName is a string containing the filename of the test
// that are being executed.
func (c *TestContext) FileName() string {
	return c.fileName
}
