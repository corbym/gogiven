package gogiven

import (
	"github.com/corbym/gogiven/base"
	"github.com/fatih/camelcase"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
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
	return some.When(action...)
}

func testTitle(functionName string) string {
	lastDotInTestName := strings.LastIndex(functionName, ".Test") + (len(".Test") - 1)
	return strings.Replace(strings.Join(camelcase.Split(functionName[lastDotInTestName+1:]), " "), "_", " ", -1)
}

// ParseGivenWhenThen parses a test file for the Given/When/Then content of the test in question identified by the parameter "testName".
// Returns the content of the function with all metacharacters removed, spaces added to CamelCase and snake case too.
func ParseGivenWhenThen(functionName string, testFileName string) (formattedOutput []string) {
	split := strings.Split(functionName, ".")
	functionName = rawFuncName(functionName, split)
	source, _ := ioutil.ReadFile(testFileName)
	fset, testFunction, _ := parseFile(testFileName, functionName, source)

	givenWhenIndex := positionOfGivenOrWhen(testFunction.Body, fset)
	source = source[givenWhenIndex:fset.Position(testFunction.End()).Offset]
	interestingStatements := removeAllUninterestingStatements(string(source[:]))
	interestingStatements = markupComments(interestingStatements)
	splitByLines := strings.Split(interestingStatements, "\n")
	formattedOutput = cleanUpGivenWhenThenOutput(splitByLines)
	return
}

func markupComments(source string) (replaced string) {
	r := regexp.MustCompile(`(/\*([^*]|[\r\n]|(\*+([^*/]|[\r\n])))*\*+/)|(//.*)`)
	replaced = r.ReplaceAllStringFunc(source, func(comment string) string {
		return "Noting that " + strings.TrimSpace(strings.Replace(comment, "/", "", -1))
	})
	return
}

func cleanUpGivenWhenThenOutput(splitByLines []string) (formattedOutput []string) {
	for _, str := range splitByLines {
		str = strings.TrimSpace(strings.Join(camelcase.Split(str), " "))
		str = strings.Replace(str, "return", "", -1)
		str = strings.TrimSpace(replaceAllNonAlphaNumericCharacters(str))
		if str != "" {
			str = strings.ToUpper(str[:1]) + strings.ToLower(str[1:])
			formattedOutput = append(formattedOutput, str)
		}
	}
	return
}

func positionOfGivenOrWhen(currentFuncBody *ast.BlockStmt, fset *token.FileSet) int {
	visitor := &identVisitor{fset: fset, fileOffsetPos: -1}
	ast.Walk(visitor, currentFuncBody)
	if visitor.fileOffsetPos == -1 {
		panic("could not find position of first given or when statement in func body")
	}
	return visitor.fileOffsetPos
}

func rawFuncName(functionName string, split []string) string {
	functionName = split[len(split)-1]
	for i := 2; functionName == "func1"; i++ {
		functionName = split[len(split)-i]
	}
	return functionName
}

func parseFile(fileName string, functionName string, src []byte) (fset *token.FileSet, fun *ast.FuncDecl, error error) {
	fset = token.NewFileSet()
	if file, err := parser.ParseFile(fset, fileName, src, parser.ParseComments); err == nil {
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
	removed = removeDeclarations(content, "interface", "{", "}")
	removed = removeDeclarations(removed, "func", "(", ")")
	return
}
func removeDeclarations(content string, keyWord string, openBracket string, closeBracket string) string {
	for strings.Contains(content, keyWord) {
		firstInstance := strings.Index(content, keyWord)
		lastBracketOfFunc := firstInstance + findBalancedBracketFor(content[firstInstance:], openBracket, closeBracket)
		funcString := content[firstInstance:lastBracketOfFunc]
		content = strings.Replace(content, funcString, "", 1)
	}
	return content
}

func findBalancedBracketFor(remove string, openBracket string, closeBracket string) (currentPosition int) {
	currentPosition = 0
	balance := -1
	for balance != 0 && currentPosition < (len(remove)-1) {
		if remove[currentPosition:currentPosition+1] == openBracket {
			if balance == -1 {
				balance++
			}
			balance++
		}
		if remove[currentPosition:currentPosition+1] == closeBracket {
			balance--
		}
		currentPosition++
	}
	return
}

func replaceAllNonAlphaNumericCharacters(replace string) (replaced string) {
	r := regexp.MustCompile("(?sm:([^a-zA-Z0-9*!£$%+/\\-^\"= \\r\\n\\t<>]))")
	replaced = r.ReplaceAllString(replace, "")
	r = regexp.MustCompile("\\s+")
	replaced = r.ReplaceAllString(replaced, " ")
	r = regexp.MustCompile(`".*?"`)
	replaced = r.ReplaceAllStringFunc(replaced, func(quoted string) string {
		return "\"" + strings.TrimSpace(strings.Replace(quoted, "\"", "", -1)) + "\""
	})
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
