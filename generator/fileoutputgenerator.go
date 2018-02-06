package generator

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var _ = mime.AddExtensionType(".json", "application/json") // default .json to application/json type as mime does not know about it

//FileOutputGenerator is a struct which implements the OutputListener interface
type FileOutputGenerator struct {
	OutputListener
}

//Notify is called by Generator to pass the output content from the test. Output is an io.Reader,
// usually test output. The file is written to env var GOGIVENS_OUTPUT_DIR if set, or defaults either
// to system tmp or the current dir if neither are found.
func (f *FileOutputGenerator) Notify(testFilePath string, contentType string, output io.Reader) {
	outputFileName := fmt.Sprintf("%s%c%s", OutputDirectory(),
		os.PathSeparator,
		f.Ref(testFilePath, contentType),
	)

	out := withResultErrorHandler(ioutil.ReadAll(output)).([]byte)

	errorHandler(ioutil.WriteFile(outputFileName, out, 0644))

	fmt.Printf("\ngenerated test output: file://%s\n", strings.Replace(outputFileName, "\\", "/", -1))
}

func (f *FileOutputGenerator) Ref(testFilePath string, contentType string) string {
	extensions := withResultErrorHandler(mime.ExtensionsByType(contentType)).([]string)

	sort.Sort(sort.Reverse(sort.StringSlice(extensions)))
	fileExtension := extensions[0]
	filename := strings.Replace(filepath.Base(testFilePath), ".go", fileExtension, 1)
	if !strings.HasSuffix(filename, fileExtension) {
		filename += fileExtension
	}
	return filename
}

func withResultErrorHandler(in interface{}, err error) interface{} {
	if err != nil {
		panic("error generating output:" + err.Error())
	}
	return in
}

func errorHandler(err error) {
	if err != nil {
		panic("error generating output:" + err.Error())
	}
}

//OutputDirectory finds the output dir in either one of the system env var $GOGIVENS_OUTPUT_DIR, the os tmp dir,
// or default to the current dir otherwise.
func OutputDirectory() string {
	outputDir := os.Getenv("GOGIVENS_OUTPUT_DIR")
	if outputDir == "" {
		os.Stdout.WriteString("env var GOGIVENS_OUTPUT_DIR was not found, using TempDir " + os.TempDir())
		outputDir = os.TempDir()
	}
	if _, err := os.Stat(outputDir); err == nil {
		return outputDir
	}
	os.Stderr.WriteString("output dir not found:" + outputDir + ", defaulting to ./")
	return "."
}
