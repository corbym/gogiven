package gogiven_test

import (
	"bytes"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven"
	"github.com/corbym/gogiven/generator"
	"io"
	"strings"
	"testing"
	"time"
)

type stubOutputListener struct {
	generator.OutputListener
	testFilePath string
	contentType  string
	output       string
	received     chan bool
}

func (stub *stubOutputListener) Notify(testFilePath string, contentType string, output io.Reader) {
	if strings.Contains(testFilePath, "generator_test") {
		stub.testFilePath = testFilePath
		stub.contentType = contentType
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(output)
		stub.output = buffer.String()
		stub.received <- true
	}
}

func newStubListener() (outputListener *stubOutputListener, hasReceived chan bool) {
	hasReceived = make(chan bool, 1)
	outputListener = &stubOutputListener{received: hasReceived}
	gogiven.OutputListeners = []generator.OutputListener{outputListener}
	defer time.AfterFunc(500*time.Millisecond, func() { hasReceived <- false })
	return
}

func TestGenerateTestOutput_contentType(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)

	listener, received := newStubListener()
	gogiven.GenerateTestOutput()
	done := <-received

	then.AssertThat(t, done, is.EqualTo(true))
	then.AssertThat(t, listener.contentType, is.EqualTo("text/html"))
}

func TestGenerateTestOutput_fileName(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)

	listener, channel := newStubListener()
	gogiven.GenerateTestOutput()
	done := <-channel

	then.AssertThat(t, done, is.EqualTo(true))
	then.AssertThat(t, listener.testFilePath, is.ValueContaining("generator_test.go"))
}

func TestGenerateTestOutput_output(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)

	listener, received := newStubListener()
	gogiven.GenerateTestOutput()
	done := <-received

	then.AssertThat(t, done, is.EqualTo(true))
	then.AssertThat(t, listener.output, is.ValueContaining("foo"))
}
