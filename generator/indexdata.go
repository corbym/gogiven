package generator

//IndexData holds information about a test for indexing.
type IndexData struct {
	// Ref is reference to where the content was deposited, produced by the OutputGenerator. E.g. /foo/bin/goo_test.html for the FileOutputGenerator
	Ref string
	// Title is the title of the test e.g. Goo Test (using the filename of the test it was generated from)
	Title string
	// TestData holds an array of TestData (which holds whether the test passed or failed)
	TestData []TestData
}

//NewIndexData creates a new IndexDat object from a generated Ref (e.g. /foo/bin/goo_test.html) and PageData
func NewIndexData(generatedRef string, pageData PageData) IndexData {
	return IndexData{
		Ref:      generatedRef,
		Title:    pageData.Title,
		TestData: pageData.TestData,
	}
}
