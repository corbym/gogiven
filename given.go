package gogiven

import (
	"regexp"
	"runtime"
	"strings"
)

var globalTestContextMap = newSafeMap()

//Given sets up some interesting givens for the test.
//Pass in testing.T here and a function which adds some givens to the map.
func Given(testing TestingT, given ...func(givens *InterestingGivens)) *Some {
	function, testFileName := testFunctionFileName()
	var currentTestContext *TestContext

	if value, ok := globalTestContextMap.Load(testFileName); ok {
		currentTestContext = value.(*TestContext)
	} else {
		currentTestContext = NewTestContext(testFileName)
		globalTestContextMap.Store(testFileName, currentTestContext)
	}
	someTests := currentTestContext.someTests
	keyFor := uniqueKeyFor(someTests, function.Name()) // this deals with table test for loops, we want different id for each

	some := NewSome(
		testing,
		testTitle(function.Name()),
		NewTestMetaData(keyFor),
		ParseGivenWhenThen(function.Name(), currentTestContext.fileContent),
		given...,
	)
	someTests.Store(keyFor, some)

	return some
}
func testTitle(functionName string) string {
	lastDotInTestName := strings.LastIndex(functionName, ".Test")
	return strings.ToTitle(functionName[lastDotInTestName+1:])
}

// ParseGivenWhenThen parses a test file for the Given/When/Then content of the test in question identified by the parameter "testName"
// Returns the content of the function with all metacharacters removed, spaces added to CamelCase and snake case too.
func ParseGivenWhenThen(name string, testFileContent string) string {
	lastDotInTestName := strings.LastIndex(name, ".")
	testName := name[lastDotInTestName+1:]

	openingBracketIndex, funcEndIndex := findOpenBracketFuncEndBracketIndices(testFileContent, testName)
	var replace = testFileContent[openingBracketIndex:funcEndIndex]
	replace = removeAllInnerFuncs(replace)
	replace = replaceAllNonAlphaNumericCharactersWithSpaces(replace)
	return replace
}

func removeAllInnerFuncs(content string) string {
	regex := regexp.MustCompile("(?m:func\\s?\\(.*\\)\\s?{)")
	content = regex.ReplaceAllString(content, " ")
	return content
}

func replaceAllNonAlphaNumericCharactersWithSpaces(replace string) string {
	r := regexp.MustCompile("(?sm:([^a-zA-Z0-9*\n\t<>]))")
	replace = r.ReplaceAllString(replace, " ")
	return replace
}

func findOpenBracketFuncEndBracketIndices(content string, funcEnd string) (openBracketIndex int, funcEndIndex int) {
	index := strings.Index(content, "func "+funcEnd)
	if index == -1{
		index = strings.Index(content, "func"+funcEnd)
	}

	funcIndex := index + len(funcEnd)
	openingBracketIndex := funcIndex + strings.Index(content[funcIndex:], "{")
	funcEndIndex = openingBracketIndex + indexForFuncEnd(content[openingBracketIndex:])
	return openingBracketIndex + 1, funcEndIndex
}

func indexForFuncEnd(findIn string) int {
	bytes := findIn[:]
	var bracketCount = -1
	for pos, char := range bytes {
		switch char {
		case '{':
			bracketCount = incBracketCountIfNegative(bracketCount)
			bracketCount++
		case '}':
			bracketCount = incBracketCountIfNegative(bracketCount)
			bracketCount--
		}
		if bracketCount == 0 {
			return pos
		}
	}
	return -1
}
func incBracketCountIfNegative(bracketCount int) int {
	if bracketCount == -1 {
		bracketCount++
	}
	return bracketCount
}

func testFunctionFileName() (*runtime.Func, string) {
	funcProgramCounters, function := findTestFpcFunction()
	testFileName, _ := function.FileLine(funcProgramCounters[0] - 1)
	return function, testFileName
}

func findTestFpcFunction() ([]uintptr, *runtime.Func) {
	funcProgramCounters := make([]uintptr, 1)
	var function *runtime.Func
	var cnt = 1
	for notFound := true; notFound; notFound = !strings.Contains(function.Name(), ".Test") {
		noOfEntries := runtime.Callers(cnt, funcProgramCounters)
		if noOfEntries == 0 {
			panic("eek, no entries in callers list; cannot set funcProgramCounters")
		}
		// get the info of the actual function that's in the pointer
		function = runtime.FuncForPC(funcProgramCounters[0] - 1)
		if function == nil {
			panic("arrgh: no function found, or dropped off end of stack!")
		}
		cnt++
	}
	return funcProgramCounters, function
}

func uniqueKeyFor(somes *SafeMap, name string) string {
	if _, ok := somes.Load(name); !ok {
		return name
	}
	return uniqueKeyFor(somes, name+"_1")
}
