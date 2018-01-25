package gogiven_test

import (
	"bytes"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven"
	"github.com/corbym/gogiven/generator"
	"io"
	"testing"
)

type stubOutputListener struct {
	generator.OutputListener
	testFilePath string
	contentType  string
	output       string
}

func (stub *stubOutputListener) Notify(testFilePath string, contentType string, output io.Reader) {
	stub.testFilePath = testFilePath
	stub.contentType = contentType
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(output)
	stub.output = buffer.String()
}

func TestGenerateTestOutput_contentType(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)

	listener := &stubOutputListener{}
	gogiven.OutputListeners = []generator.OutputListener{listener}
	gogiven.GenerateTestOutput()
	then.AssertThat(t, listener.contentType, is.EqualTo("text/html"))
}

func TestGenerateTestOutput_fileName(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)

	listener := &stubOutputListener{}
	gogiven.OutputListeners = []generator.OutputListener{listener}
	gogiven.GenerateTestOutput()
	then.AssertThat(t, listener.testFilePath, is.ValueContaining("generator_test.go"))
}

func TestGenerateTestOutput_output(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)

	listener := &stubOutputListener{}
	gogiven.OutputListeners = []generator.OutputListener{listener}
	gogiven.GenerateTestOutput()
	then.AssertThat(t, listener.output, is.ValueContaining("foo"))
}
