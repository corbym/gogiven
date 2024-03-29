package generator

import (
	"bytes"
	"fmt"
	"github.com/corbym/gocrest"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"runtime"
	"sort"
	"strings"
	"testing"
)

const theContent = "generated content"

var underTest = &FileOutputGenerator{}

func TestFileOutputGenerator_Notify(t *testing.T) {
	tests := []struct {
		name                     string
		contentType              string
		fullyQualifiedFileName   string
		fileNameWithoutExtension string
	}{
		{"json content", "application/json", "./funbags.go", "funbags"},
		{"html content", "text/html", "./funbags.go", "funbags"},
		{"pdf content", "application/pdf", "./funbags.go", "funbags"},
		{"pdf content", "application/pdf", "index", "index"},
	}
	for _, testRange := range tests {
		t.Run(testRange.name, func(t *testing.T) {
			// set up
			extension := firstOfSortedExtensions(testRange.contentType)
			expectedFileOutputFileName := testRange.fileNameWithoutExtension + extension[0]

			defer func() {
				then.AssertThat(t, someFile(ofFileInOutputDir(expectedFileOutputFileName)), exists())
				then.AssertThat(t, theContentOfThe(expectedFileOutputFileName), is.EqualTo(theContent))
				removeFileName := ofFileInOutputDir(expectedFileOutputFileName)
				remove := os.Remove(removeFileName)
				then.AssertThat(t, remove, is.Nil())
			}()

			reader := strings.NewReader(theContent)
			underTest.Notify(testRange.fullyQualifiedFileName, testRange.contentType, reader)

		})
	}
}
func TestFileOutputGenerator_panics_IncorrectContentType(t *testing.T) {
	defer func() {
		panics := recover()
		then.AssertThat(t, panics, is.Not(is.Nil()))
	}()

	underTest.Notify("./flap.foo", "widget/fong", strings.NewReader(theContent))
}

func TestFileOutputGenerator_panics_ReadingContent(t *testing.T) {
	defer func() {
		panics := recover()
		then.AssertThat(t, panics, is.Not(is.Nil()))
	}()
	underTest.Notify("./flap.foo", "text/html", &mockErroringReader{})
}

func TestFileOutputGenerator_panics_WritingFile(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.SkipNow()
	}
	defer func() {
		panics := recover()
		then.AssertThat(t, panics, is.Not(is.Nil()))
	}()

	underTest.Notify("./f****0**.go", "text/html", strings.NewReader(theContent))
}

func TestGenerateTestOutput_DefaultsToCurrentDir(t *testing.T) {
	old := os.Getenv("GOGIVENS_OUTPUT_DIR")
	extension := firstOfSortedExtensions("text/html")
	expectedFileOutputFileName := "funbags" + extension[0]

	defer func() {
		then.AssertThat(t, someFile("./"+expectedFileOutputFileName), exists())
		removeFileName := "./" + expectedFileOutputFileName
		remove := os.Remove(removeFileName)
		then.AssertThat(t, remove, is.Nil())
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
	content, _ := ioutil.ReadFile(ofFileInOutputDir(expectedFileOutput))
	return string(content[:])
}

func someFile(pathToFile string) os.FileInfo {
	fileInfo, err := os.Stat(pathToFile)
	if err != nil {
		panic(err)
	}
	return fileInfo
}

func exists() *gocrest.Matcher {
	matcher := new(gocrest.Matcher)
	matcher.Matches = func(actual interface{}) bool {
		file, ok := actual.(os.FileInfo)
		if ok {
			matcher.Describe = file.Name()
			return true
		}
		return false
	}
	return matcher
}

func ofFileInOutputDir(fileName string) string {
	return fmt.Sprintf("%s%c%s", OutputDirectory(), os.PathSeparator, fileName)
}

type mockErroringReader struct {
	io.Reader
}

func (*mockErroringReader) Read(p []byte) (n int, err error) {
	err = bytes.ErrTooLarge // ReadAll only errors when this is the error from the io.Read method
	return
}
