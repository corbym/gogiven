package gogiven

import (
	"regexp"
	"runtime"
	"strings"
	"go/parser"
	"go/token"
	"go/ast"
	bytes2 "bytes"
	"go/format"
	"github.com/fatih/camelcase"
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
	lastDotInTestName := strings.LastIndex(functionName, ".Test") + (len(".Test") - 1)
	return strings.Join(camelcase.Split(functionName[lastDotInTestName+1:]), " ")
}

// ParseGivenWhenThen parses a test file for the Given/When/Then content of the test in question identified by the parameter "testName"
// Returns the content of the function with all metacharacters removed, spaces added to CamelCase and snake case too.
func ParseGivenWhenThen(functionName string, testFileContent string) string {
	fset := token.NewFileSet()
	testFile, err := parser.ParseFile(fset, functionName, testFileContent, 0)
	if err != nil {
		panic(err.Error())
	}
	line := make([]byte, 0)
	buffer := bytes2.NewBuffer(line)
	for _, dcl := range testFile.Decls {
		if fn, ok := dcl.(*ast.FuncDecl); ok {
			if strings.Contains(functionName, fn.Name.Name) {
				for _, statement := range fn.Body.List {
					if exprStmt, ok := statement.(*ast.ExprStmt); ok {
						if call, ok := exprStmt.X.(*ast.CallExpr); ok {
							if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
								funcName := fun.Sel.Name
								if funcName == "Given" || funcName == "When" || funcName == "Then" {
									format.Node(buffer, fset, statement)
								}
							}
						}
					}
				}

			}
		}
	}

	return strings.TrimSpace(replaceAllNonAlphaNumericCharactersWithSpaces(
		removeAllInnerFuncs(buffer.String()),
	))
}

func removeAllInnerFuncs(content string) string {
	regex := regexp.MustCompile("(?m:func\\s?\\(.*\\)\\s?{)")
	content = regex.ReplaceAllString(content, " ")
	return content
}

func replaceAllNonAlphaNumericCharactersWithSpaces(replace string) string {
	r := regexp.MustCompile("(?sm:([^a-zA-Z0-9*\"\n\t<>]))")
	replace = r.ReplaceAllString(replace, " ")
	return replace
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
