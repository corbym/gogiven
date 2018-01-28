package gogiven_test

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven"
	"github.com/corbym/gogiven/test_stubs"
	"testing"
)

func TestGenerateTestOutput_contentType(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)
	listener, received := test_stubs_test.NewStubListener()
	gogiven.GenerateTestOutput()
	done := <-received

	then.AssertThat(t, done, is.EqualTo(true))
	then.AssertThat(t, listener.ContentType, is.EqualTo("text/html"))
}

func TestGenerateTestOutput_fileName(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)

	listener, channel := test_stubs_test.NewStubListener()
	gogiven.GenerateTestOutput()
	done := <-channel

	then.AssertThat(t, done, is.EqualTo(true))
	then.AssertThat(t, listener.TestFilePath, is.ValueContaining("generator_test.go"))
}

func TestGenerateTestOutput_output(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)

	listener, received := test_stubs_test.NewStubListener()
	gogiven.GenerateTestOutput()
	done := <-received

	then.AssertThat(t, done, is.EqualTo(true))
	then.AssertThat(t, listener.Output, is.ValueContaining("foo"))
}

func TestGenerateTestOutput_GenerateIndex(t *testing.T) {
	oldOutputGenerator := gogiven.Generator
	defer func() {
		gogiven.Generator = oldOutputGenerator
	}()
	gogiven.Given(t)
	_, received := test_stubs_test.NewStubGenerator()
	gogiven.GenerateTestOutput()
	done := <-received

	then.AssertThat(t, done, is.EqualTo(true))
}
