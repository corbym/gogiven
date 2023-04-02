package generator

import (
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/testdata"
)

// TestResult holds the results of the test, whether it failed, was skipped, it's ID and
// the output from the Fail/Skip if appropriate
type TestResult struct {
	TestID     string
	Failed     bool
	Skipped    bool
	TestOutput string
}

// NewTestResult returns a copy of the test meta data
func NewTestResult(data *base.TestMetaData) TestResult {
	return TestResult{
		Skipped:    data.Skipped(),
		TestID:     data.Name(),
		TestOutput: data.TestOutput(),
		Failed:     data.Failed(),
	}
}

// TestData holds a copy of test results, interesting givens, captured io, and parsed test content (given/when/then)
type TestData struct {
	TestResult        TestResult
	TestTitle         string
	InterestingGivens testdata.InterestingGivens
	CapturedIO        testdata.CapturedIO
	ParsedTestContent base.ParsedTestContent
}

// NewTestData creates a new copy of test results, interesting givens, captured io, and parsed test content (given/when/then)
func NewTestData(some *base.Some) (testData TestData) {
	testData = TestData{
		TestResult:        NewTestResult(some.TestMetaData()),
		InterestingGivens: some.InterestingGivens(),
		CapturedIO:        some.CapturedIO(),
		ParsedTestContent: some.ParsedTestContent(),
		TestTitle:         some.TestTitle(),
	}
	return
}

// PageData is the struct that populates the template with data from the test output.
type PageData struct {
	Title    string
	TestData []TestData
}

// NewPageData holds the title and the test results
func NewPageData(title string, someMap *base.SomeMap) (pageData PageData) {
	pageData = PageData{
		Title:    title,
		TestData: copyTestResults(someMap),
	}
	return
}

func copyTestResults(someMap *base.SomeMap) []TestData {
	var testData []TestData

	for _, v := range *someMap {
		testData = append(testData, NewTestData(v))
	}
	return testData
}
