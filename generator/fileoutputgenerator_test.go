package generator

import (
	"fmt"
	"github.com/corbym/gocrest"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type args struct {
	contentType string
	expectedExt string
}

func TestFileOutputGenerator_Notify(t *testing.T) {
	const theContent = "generated content"
	underTest := &FileOutputGenerator{}

	tests := []struct {
		name string
		args args
	}{
		{"json content", args{contentType: "application/json", expectedExt: ".json"}},
		{"html content", args{contentType: "text/html", expectedExt: ".html"}},
	}
	for _, testRange := range tests {
		t.Run(testRange.name, func(t *testing.T) {
			// set up
			expectedFileOutputFileName := "generator_test" + testRange.args.expectedExt
			withContentType := testRange.args.contentType
			defer func() {
				then.AssertThat(t, someFileExists(ofFileInTmpDir(expectedFileOutputFileName)), inTmpDir())
				then.AssertThat(t, theContentOfThe(expectedFileOutputFileName), is.EqualTo(theContent))
				remove := os.Remove(ofFileInTmpDir(expectedFileOutputFileName))
				then.AssertThat(t, remove, is.Nil())
			}()

			reader := strings.NewReader(theContent)
			underTest.Notify("./generator_test.go", withContentType, reader)

		})
	}
}

func TestGenerateTestOutput_DefaultsToCurrentDir(t *testing.T) {
	old := os.Getenv("GOGIVENS_OUTPUT_DIR")
	defer func() {
		then.AssertThat(t, someFileExists("./generator_test.html"), inTmpDir())
		os.Remove("./generator_test.html")
	}()
	defer func() { os.Setenv("GOGIVENS_OUTPUT_DIR", old) }()
	os.Setenv("GOGIVENS_OUTPUT_DIR", "doesnotexist")
	outputGenerator := &FileOutputGenerator{}

	outputGenerator.Notify("./generator_test.go", "text/html", strings.NewReader("foo"))

}

func theContentOfThe(expectedFileOutput string) string {
	content, _ := ioutil.ReadFile(ofFileInTmpDir(expectedFileOutput))
	return string(content[:])
}

func someFileExists(pathToFile string) interface{} {
	fileInfo, err := os.Stat(pathToFile)
	if err != nil {
		return err
	}
	return fileInfo
}

func inTmpDir() *gocrest.Matcher {
	matcher := new(gocrest.Matcher)
	matcher.Matches = func(actual interface{}) bool {
		file, ok := actual.(os.FileInfo)
		if ok {
			matcher.Describe = fmt.Sprintf("%s", file.Name())
			return true
		}
		return false
	}
	return matcher
}
func ofFileInTmpDir(fileName string) string {
	return fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, fileName)
}
