package gogiven

import (
	"fmt"
	"io/ioutil"
	"strings"
	"os"
	"path/filepath"
)

type HtmlGenerator interface {
	Generate(fileNameWithPath string, testFileContent string) (html string)
}

type HtmlFileGenerator struct {
	HtmlGenerator *HtmlFileGenerator
}

func (generator *HtmlFileGenerator) Generate(fileNameWithPath string, testFileContent string) (html string) {
	return "Ok"
}

var Generator HtmlGenerator = new(HtmlFileGenerator)

func generateTestOutput(some *some) {
	testFileName, content := readTestFile(some)
	html := Generator.Generate(testFileName, string(content[:]))
	fileName := fmt.Sprintf("%s%c%s", os.TempDir(),
		os.PathSeparator,
		strings.Replace(filepath.Base(testFileName), ".go", ".html", 1))

	errWritingHtml := ioutil.WriteFile(fileName, []byte(html), 0644)
	if errWritingHtml != nil {
		panic("error writing html:" + errWritingHtml.Error())
	}
	fmt.Printf("generated test output: file:///%s\n", strings.Replace(fileName, "\\", "/", -1))
}

func readTestFile(some *some) (string, []byte) {
	file, _ := some.runtimeCaller.FileLine(some.frameProgramCounters[0] - 1) // second result is "line" number
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic("file not found:" + err.Error())
	}
	return file, content
}
