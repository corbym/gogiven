package generator

import (
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/testdata"
)

type TestResult struct {
	TestId     string
	Failed     bool
	Skipped    bool
	TestOutput string
}

func NewTestResult(data *base.TestMetaData) (result *TestResult) {
	result = &TestResult{
		Skipped:    data.Skipped(),
		TestId:     data.Name(),
		TestOutput: data.TestOutput(),
		Failed:     data.Failed(),
	}
	return
}

type TestData struct {
	TestResult        *TestResult
	TestTitle         string
	InterestingGivens testdata.InterestingGivens
	CapturedIO        testdata.CapturedIO
	ParsedTestContent base.ParsedTestContent
}

func NewTestData(some *base.Some) (testData *TestData) {
	testData = &TestData{
		TestResult:        NewTestResult(some.TestMetaData()),
		InterestingGivens: some.InterestingGivens(),
		CapturedIO:        some.CapturedIO(),
		ParsedTestContent: some.GivenWhenThen(),
		TestTitle:         some.TestTitle(),
	}
	return
}

//PageData is the struct that populates the template with data from the test output.
type PageData struct {
	Title       string
	TestResults map[string]*TestData
}

func NewPageData(title string, someMap *base.SomeMap) (pageData *PageData) {
	pageData = new(PageData)
	pageData.Title = title
	pageData.TestResults = copyTestResults(someMap)
	return
}

func copyTestResults(someMap *base.SomeMap) (testData map[string]*TestData) {
	testData = make(map[string]*TestData)
	for k, v := range *someMap {
		testData[k] = NewTestData(v)
	}
	return
}
