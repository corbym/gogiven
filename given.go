package gogiven

import (
	"github.com/corbym/gogiven/base"
	"github.com/fatih/camelcase"
	"runtime"
	"strings"
)

var globalTestContextMap = newSafeMap()

//Given sets up some interesting givens for the test.
//Pass in testing.T here and a function which adds some givens to
//the map.
//
// *Warning:* if you call this method twice in a test, you will start a new test
//which will appear in the test output.
func Given(testing base.TestingT, given ...base.GivenData) *base.Some {
	currentFunction, testFileName := testFunctionFileName()
	currentTestContext := loadTestContext(testFileName)

	someTests := currentTestContext.someTests
	keyFor := uniqueKeyFor(someTests, currentFunction.Name()) // this deals with table test for loops, we want different id for each

	some := base.NewSome(
		testing,
		testTitle(testing.Name()),
		base.NewTestMetaData(keyFor),
		base.ParseGivenWhenThen(currentFunction.Name(), currentTestContext.FileName()),
		given...,
	)
	someTests.Store(keyFor, some)

	return some
}

func loadTestContext(testFileName string) (currentTestContext *TestContext) {
	if value, ok := globalTestContextMap.Load(testFileName); ok {
		currentTestContext = value.(*TestContext)
	} else {
		currentTestContext = NewTestContext(testFileName)
		globalTestContextMap.Store(testFileName, currentTestContext)
	}
	return
}

//When is a shortcut method when no Given is required.
func When(testing base.TestingT, action ...base.CapturedIOGivenData) *base.Some {
	some := Given(testing)
	return some.When(action...)
}

func testTitle(functionName string) string {
	lastDotInTestName := strings.LastIndex(functionName, ".Test") + (len(".Test") - 1)
	return strings.Replace(strings.Join(camelcase.Split(functionName[lastDotInTestName+1:]), " "), "_", " ", -1)
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
