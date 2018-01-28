package test_stubs_test

import (
	"bytes"
	"github.com/corbym/gogiven"
	"github.com/corbym/gogiven/generator"
	"io"
	"strings"
	"time"
)

type StubOutputListener struct {
	generator.OutputListener
	TestFilePath string
	ContentType  string
	Output       string
	Received     chan bool
}

func NewStubListener() (outputListener *StubOutputListener, hasReceived chan bool) {
	hasReceived = make(chan bool, 1)
	outputListener = &StubOutputListener{Received: hasReceived}
	gogiven.OutputListeners = []generator.OutputListener{outputListener}
	defer time.AfterFunc(500*time.Millisecond, func() { hasReceived <- false })
	return
}

func (stub *StubOutputListener) Notify(testFilePath string, contentType string, output io.Reader) {
	if strings.Contains(testFilePath, "generator_test") {
		stub.TestFilePath = testFilePath
		stub.ContentType = contentType
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(output)
		stub.Output = buffer.String()
		stub.Received <- true
	}
}

type StubGenerator struct {
	generator.GoGivensOutputGenerator
	Received chan bool
}

func NewStubGenerator() (generator *StubGenerator, hasReceived chan bool) {
	hasReceived = make(chan bool, 1)
	generator = &StubGenerator{Received: hasReceived}
	gogiven.Generator = generator
	defer time.AfterFunc(500*time.Millisecond, func() { hasReceived <- false })
	return
}
func (stubGenerator *StubGenerator) Generate(data *generator.PageData) (output io.Reader) {
	return strings.NewReader("content")
}

func (stubGenerator *StubGenerator) GenerateIndex() (output io.Reader) {
	stubGenerator.Received <- true
	return strings.NewReader("index")
}

func (stubGenerator *StubGenerator) ContentType() string {
	return "text/html"
}
