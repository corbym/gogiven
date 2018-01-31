package generator

import (
	"io"
)

//OutputListener is the interface of things that want to be notified about generated content.
type OutputListener interface {
	Notify(testFilePath string, contentType string, output io.Reader)
}
