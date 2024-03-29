package gogiven_test

import (
	"github.com/corbym/gocrest/has"
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
	listener, received := test_stubs.NewStubListener(test_stubs.GeneratorTest)
	gogiven.GenerateTestOutput()
	done := <-received

	then.AssertThat(t, done, is.EqualTo(test_stubs.GeneratorTest))
	then.AssertThat(t, listener.ContentType, is.EqualTo("text/html"))
}

func TestGenerateTestOutput_fileName(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)
	listener, received := test_stubs.NewStubListener(test_stubs.GeneratorTest)
	gogiven.GenerateTestOutput()
	done := <-received

	then.AssertThat(t, done, is.EqualTo(test_stubs.GeneratorTest))
	then.AssertThat(t, listener.TestFilePath, is.ValueContaining("generator_test.go"))
}

func TestGenerateTestOutput_output(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	defer func() {
		gogiven.OutputListeners = oldListeners
	}()
	// initialise global map
	gogiven.Given(t)

	listener, received := test_stubs.NewStubListener(test_stubs.GeneratorTest)
	gogiven.GenerateTestOutput()
	<-received

	then.AssertThat(t, listener.Output, is.ValueContaining("foo"))
}

func TestGenerateTestOutput_GenerateIndex(t *testing.T) {
	oldOutputGenerator := gogiven.Generator
	defer func() {
		gogiven.Generator = oldOutputGenerator
	}()

	gogiven.Given(t)
	stubGenerator, received := test_stubs.NewStubGenerator()

	gogiven.GenerateTestOutput()
	done := <-received

	then.AssertThat(t, done, is.EqualTo(true))
	then.AssertThat(t, stubGenerator.IndexData, has.Length(is.GreaterThan(0)))
	then.AssertThat(t, stubGenerator.IndexData[0].Ref, is.AnyOf(
		is.ValueContaining("_test.html"),
		is.ValueContaining("_test.shtml"),
	))
}

func TestGenerateTestOutput_OutputIndex(t *testing.T) {
	oldListeners := gogiven.OutputListeners
	oldOutputGenerator := gogiven.Generator
	defer func() {
		gogiven.OutputListeners = oldListeners
		gogiven.Generator = oldOutputGenerator
	}()

	// initialise global map
	gogiven.Given(t)

	listener, received := test_stubs.NewStubListener(test_stubs.IndexFileName)
	test_stubs.NewStubGenerator()

	gogiven.GenerateTestOutput()
	done := <-received

	then.AssertThat(t, done, is.EqualTo(test_stubs.IndexFileName))
	then.AssertThat(t, listener.Output, is.EqualTo("index"))
}
