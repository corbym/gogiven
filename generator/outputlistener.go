package generator

import (
	"io"
)

//OutputListener is the interface of things that want to be notified about generated content.
type OutputListener interface {
	//Ref should generate a reference to the content, if relevant for the listener
	Ref(testFilePath string, contentType string) string
	//Notify is called by GoGiven when content is generated.
	Notify(testFilePath string, contentType string, output io.Reader)
}
