package generator

import (
	"io"
)

type OutputListener interface {
	Notify(testFilePath string, contentType string, output io.Reader)
}
