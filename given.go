package gogiven

import (
	"bytes"
	"github.com/corbym/gogiven/base"
	"github.com/fatih/camelcase"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
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
		ParseGivenWhenThen(currentFunction.Name(), currentTestContext.fileName),
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
	return some.When(action ...)
}

func testTitle(functionName string) string {
	lastDotInTestName := strings.LastIndex(functionName, ".Test") + (len(".Test") - 1)
	return strings.Replace(strings.Join(camelcase.Split(functionName[lastDotInTestName+1:]), " "), "_", " ", -1)
}

// ParseGivenWhenThen parses a test file for the Given/When/Then content of the test in question identified by the parameter "testName"
// Returns the content of the function with all metacharacters removed, spaces added to CamelCase and snake case too.
func ParseGivenWhenThen(functionName string, testFileName string) (formattedOutput string) {

	buffer := new(bytes.Buffer)

	split := strings.Split(functionName, ".")
	functionName = split[len(split)-1]
	for i := 2; functionName == "func1"; i++ {
		functionName = split[len(split)-i]
	}
	fset, fun, _ := parseFile(testFileName, functionName)
	format.Node(buffer, fset, fun.Body)

	formattedOutput = buffer.String()
	formattedOutput = formattedOutput[1: len(formattedOutput)-1]
	formattedOutput = strings.Join(camelcase.Split(removeAllUninterestingStatements(formattedOutput)), " ")
	formattedOutput = replaceAllNonAlphaNumericCharacters(formattedOutput)
	formattedOutput = strings.TrimSpace(strings.Replace(formattedOutput, "\n\t", "\n", -1))
	return
}

func parseFile(fileName string, functionName string) (fset *token.FileSet, fun *ast.FuncDecl, error error) {
	fset = token.NewFileSet()
	if file, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments); err == nil {
		for _, d := range file.Decls {
			if f, ok := d.(*ast.FuncDecl); ok && f.Name.Name == functionName {
				fun = f
				return
			}
		}
	}
	panic("could not find function " + functionName)
}

func removeAllUninterestingStatements(content string) (removed string) {
	regex := regexp.MustCompile("(?sm:func\\s?\\(.*?\\)\\s?.*?}?)")
	removed = regex.ReplaceAllString(content, "")

	index := strings.Index(removed, "Given")
	if index == -1 {
		index = strings.Index(removed, "When")
	}
	removed = removed[index:]
	return
}

func replaceAllNonAlphaNumericCharacters(replace string) (replaced string) {
	r := regexp.MustCompile("(?sm:([^a-zA-Z0-9*!£$%+\\-^\"= \\r\\n\\t<>]))")
	replace = r.ReplaceAllString(replace, "")
	r = regexp.MustCompile("  +")
	replace = r.ReplaceAllString(replace, " ")
	r = regexp.MustCompile("\t\t+")
	replace = r.ReplaceAllString(replace, "\t\t")
	r = regexp.MustCompile("[\r\n]+")
	replaced = r.ReplaceAllString(replace, "\r\n")
	return
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

func uniqueKeyFor(somes *safeMap, name string) string {
	if _, ok := somes.Load(name); !ok {
		return name
	}
	return uniqueKeyFor(somes, name+"_1")
}