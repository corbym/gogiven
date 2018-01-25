package generator

import (
	"fmt"
	"github.com/corbym/gocrest"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"io/ioutil"
	"mime"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestFileOutputGenerator_Notify(t *testing.T) {
	const theContent = "generated content"
	underTest := &FileOutputGenerator{}
	tests := []struct {
		name        string
		contentType string
	}{
		{"json content", "application/json"},
		{"html content", "text/html"},
	}
	for _, testRange := range tests {
		t.Run(testRange.name, func(t *testing.T) {
			// set up
			extension := firstOfSortedExtensions(testRange.contentType)
			expectedFileOutputFileName := "funbags" + extension[0]

			defer func() {
				then.AssertThat(t, someFileExists(ofFileInTmpDir(expectedFileOutputFileName)), inTmpDir())
				then.AssertThat(t, theContentOfThe(expectedFileOutputFileName), is.EqualTo(theContent))
				remove := os.RemoveAll(ofFileInTmpDir("./funbags" + extension[0]))
				then.AssertThat(t, remove, is.Nil())
			}()

			reader := strings.NewReader(theContent)
			underTest.Notify("./funbags.go", testRange.contentType, reader)

		})
	}
}

func TestGenerateTestOutput_DefaultsToCurrentDir(t *testing.T) {
	old := os.Getenv("GOGIVENS_OUTPUT_DIR")
	extension := firstOfSortedExtensions("text/html")
	expectedFileOutputFileName := "funbags" + extension[0]

	defer func() {
		then.AssertThat(t, someFileExists("./"+expectedFileOutputFileName), inTmpDir())
		os.Remove("./funbags.*")
	}()
	defer func() { os.Setenv("GOGIVENS_OUTPUT_DIR", old) }()
	os.Setenv("GOGIVENS_OUTPUT_DIR", "doesnotexist")
	outputGenerator := &FileOutputGenerator{}
	outputGenerator.Notify("./funbags.go", "text/html", strings.NewReader("foo"))

}
func firstOfSortedExtensions(contentType string) []string {
	extension, _ := mime.ExtensionsByType(contentType)
	sort.Sort(sort.Reverse(sort.StringSlice(extension)))
	return extension
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
