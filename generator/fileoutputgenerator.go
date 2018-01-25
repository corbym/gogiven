package generator

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

var _ = mime.AddExtensionType(".json", "application/json") // default .json to application/json type as mime does not know about it

type FileOutputGenerator struct {
	OutputListener
}

//Notify is called by Generator to pass the output content from the test. Output is an io.Reader
func (f *FileOutputGenerator) Notify(testFilePath string, contentType string, output io.Reader) {
	fileExtension := findLongestExtensionFor(contentType) // screw you, windows

	outputFileName := fmt.Sprintf("%s%c%s", outputDirectory(),
		os.PathSeparator,
		strings.Replace(filepath.Base(testFilePath), ".go", fileExtension, 1))

	out, err := ioutil.ReadAll(output)

	err = ioutil.WriteFile(outputFileName, out, 0644)
	if err != nil {
		panic("error generating gogiven output:" + err.Error())
	}
	fmt.Printf("\ngenerated test output: file://%s\n", strings.Replace(outputFileName, "\\", "/", -1))
}

func findLongestExtensionFor(contentType string) string {
	types, _ := mime.ExtensionsByType(contentType)
	longest := 0
	retIdx := -1
	for idx, extension := range types {
		if len(extension) > longest {
			longest = len(extension)
			retIdx = idx
		}
	}
	return types[retIdx]
}

func outputDirectory() string {
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
