package gogiven

import (
	"runtime"
	"strings"
)

var globalTestContextMap = newSafeMap()

func Given(testing TestingT, given ...func(givens *InterestingGivens)) *Some {
	function, testFileName := testFunctionFileName()
	var currentTestContext *TestContext

	if value, ok := globalTestContextMap.Load(testFileName); ok {
		currentTestContext = value.(*TestContext)
	} else {
		currentTestContext = newGlobalTestContext(testFileName)
		globalTestContextMap.Store(testFileName, currentTestContext)
	}
	someTests := currentTestContext.someTests
	keyFor := uniqueKeyFor(someTests, function.Name()) // this deals with table test for loops, we want different id for each

	some := newSome(testing, newTestMetaData(keyFor), given...)
	someTests.Store(keyFor, some)

	return some
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
