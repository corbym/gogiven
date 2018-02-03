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
	Received     chan string
	await        string
}

const GeneratorTest = "generator_test"
const IndexFileName = "index"
const FailedToReceive = "failed"

func NewStubListener(await string) (*StubOutputListener, chan string) {
	hasReceived := make(chan string, 1)
	outputListener := &StubOutputListener{Received: hasReceived, await: await}
	gogiven.OutputListeners = []generator.OutputListener{outputListener}
	defer time.AfterFunc(500*time.Millisecond, func() {
		hasReceived <- FailedToReceive
	})
	return outputListener, hasReceived
}

func (stub *StubOutputListener) Notify(testFilePath string, contentType string, output io.Reader) {
	if onlyReadWhenFilePathContains(testFilePath, stub.await) {
		stub.TestFilePath = testFilePath
		stub.ContentType = contentType
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(output)
		stub.Output = buffer.String()
		stub.Received <- stub.await
	}
}
func onlyReadWhenFilePathContains(testFilePath string, await string) bool {
	return strings.Contains(testFilePath, await)
}

type StubGenerator struct {
	generator.GoGivensOutputGenerator
	Received chan bool
}

func NewStubGenerator() (*StubGenerator, chan bool) {
	hasReceived := make(chan bool, 1)
	generator := &StubGenerator{Received: hasReceived}
	gogiven.Generator = generator
	defer time.AfterFunc(2*time.Second, func() { hasReceived <- false })
	return generator, hasReceived
}
func (stubGenerator *StubGenerator) Generate(data generator.PageData) (output io.Reader) {
	return strings.NewReader("content")
}

func (stubGenerator *StubGenerator) GenerateIndex() (output io.Reader) {
	stubGenerator.Received <- true
	return strings.NewReader("index")
}

func (stubGenerator *StubGenerator) ContentType() string {
	return "text/html"
}