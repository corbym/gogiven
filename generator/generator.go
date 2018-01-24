package generator

import (
	"io"
)

// GoGivensOutputGenerator is an interface that can be implemented by anything that can generate file content to be output
// after a test has completed.
type GoGivensOutputGenerator interface {
	Generate(data *PageData) (output io.Reader)
	//ContentType is text/html, application/json or other mime type
	ContentType() string
}
