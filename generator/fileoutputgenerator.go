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

type FileOutputGenerator struct {
	OutputListener
}

func (f *FileOutputGenerator) Notify(testFilePath string, contentType string, output io.Reader) {
	fileExtension, _ := mime.ExtensionsByType(contentType)
	outputFileName := fmt.Sprintf("%s%c%s", outputDirectory(),
		os.PathSeparator,
		strings.Replace(filepath.Base(testFilePath), ".go", fileExtension[0], 1))
	out, _ := ioutil.ReadAll(output)
	err := ioutil.WriteFile(outputFileName, out, 0644)
	if err != nil {
		panic("error generating gogiven output:" + err.Error())
	}
	fmt.Printf("\ngenerated test output: file://%s\n", strings.Replace(outputFileName, "\\", "/", -1))
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
