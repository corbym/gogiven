package generator

//IndexData holds information about a test for indexing.
type IndexData struct {
	// TestFileName is the original test's full path and filename, eg. /foo/bin/goo_test.go
	TestFileName string
	// Title is the title of the test e.g. Goo Test (using the filename above)
	Title string
	// TestData holds an array of TestData (which holds whether the test passed or failed)
	TestData []TestData
}

//NewIndexData creates a new IndexDat object from a file name (/foo/bin/goo_test.go) and PageData
func NewIndexData(testFileName string, pageData PageData) IndexData {
	return IndexData{
		TestFileName: testFileName,
		Title:        pageData.Title,
		TestData:     pageData.TestData,
	}
}
