package base

import (
	"github.com/fatih/camelcase"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"regexp"
	"strings"
)

// ParsedTestContent parsed content from test
type ParsedTestContent struct {
	GivenWhenThen []string
	// Comment contains each comment line from the tests' comment block
	Comment []string
}

// ParseGivenWhenThen parses a test file for the Given/When/Then content of the test in question identified by the parameter "testName".
// Returns the content of the function with all metacharacters removed, spaces added to CamelCase and snake case too.
func ParseGivenWhenThen(functionName string, testFileName string) ParsedTestContent {
	split := strings.Split(functionName, ".")
	functionName = rawFuncName(functionName, split)
	source, _ := ioutil.ReadFile(testFileName)
	fset, testFunction, _ := parseFile(testFileName, functionName, source)
	givenWhenIndex := positionOfGivenOrWhen(testFunction.Body, fset)
	source = source[givenWhenIndex:fset.Position(testFunction.End()).Offset]
	interestingStatements := removeAllUninterestingStatements(string(source[:]))
	interestingStatements = markupComments(interestingStatements)
	splitByLines := strings.Split(interestingStatements, "\n")

	return ParsedTestContent{
		GivenWhenThen: cleanUpGivenWhenThenOutput(splitByLines),
		Comment:       strings.Split(testFunction.Doc.Text(), "\n"),
	}
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

func positionOfGivenOrWhen(currentFuncBody ast.Node, fset *token.FileSet) int {
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
	r := regexp.MustCompile(`(?sm:([^a-zA-Z0-9*!£$%+/\\-^"= \r\n\t<>]))`)
	replaced = r.ReplaceAllString(replace, "")
	r = regexp.MustCompile(`\s+`)
	replaced = r.ReplaceAllString(replaced, " ")
	r = regexp.MustCompile(`".*?"`)
	replaced = r.ReplaceAllStringFunc(replaced, func(quoted string) string {
		return `"` + strings.TrimSpace(strings.Replace(quoted, `"`, "", -1)) + `"`
	})
	return
}
