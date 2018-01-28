# gogiven
An alternative BDD spec framework for go. Builds on "go test" tool and builds on the go testing package.

[![Build status](https://travis-ci.org/corbym/gogiven.svg?branch=master)](https://travis-ci.org/corbym/gogiven/builds)
[![Go Report Card](https://goreportcard.com/badge/github.com/corbym/gogiven)](https://goreportcard.com/report/github.com/corbym/gogiven)
[![GoDoc](https://godoc.org/github.com/corbym/gogiven?status.svg)](http://godoc.org/github.com/corbym/gogiven)
[![Coverage Status](https://coveralls.io/repos/github/corbym/gogiven/badge.svg?branch=master)](https://coveralls.io/github/corbym/gogiven?branch=master)

Inspired by [YATSPEC](https://github.com/bodar/yatspec). Another similar idea is [JGiven](http://jgiven.org), although both are Java based.

Check out the HTML output generator used as a default here: [htmlspec](https://github.com/corbym/htmlspec) - using go's own template/html package.

Feel free to contact me and help improve the code base or let me know if you have any issues or questions!

To contribute, please read [this first.](https://github.com/corbym/gogiven/blob/master/CONTRIBUTING.md)

## Table of Contents
1. [Introduction](#introduction)
2. [Example One - GoGivens in practice](#example)
3. [Example Two - Table Tests](#tabletest-example)
4. [Content Generation](#content-gen)
5. [List of pre-written output generators](#output-generator-list)

## Introduction <a name="introduction"></a>

Go Givens is a lightweight TDD framework for producing test specifications directly from the code you write.

Go Givens parses your test file and produces a human-readable output in a specified directory, containing all the tests, captured data and other related information regarding your test such as success or failure.

Go Givens was inspired by YATSPEC, a BDD framework employed extensively by Sky Network Services (part of Sky, a UK tv company). As mentioned above, another similar product is [JGiven](http://jgiven.org).

### Why?

Capturing your test method as test output is the only real way to show it's intention. You can refactor a test, and have the output update accordingly when the test runs. Unlike other go BDD frameworks, you can use function names to declare intent, and refactoring the function will affect the test. E.g.

```go
Given(testing, someData)..
..

func someData(..)
```

.. will be rendered as:

```
Given testing some data
```

Test data (set by the func ```someData``` in the example above) can be captured in a map of interesting givens, and as the test progresses, the actuals can be captured in a map of captured inputs and outputs.


### Interesting whats?

Interesting givens are data points that you want to use in your tests that are important for it to function correctly. "Given" data is then used by your application to fulfil some sort of request or process.

Interesting givens are generated as output along side your test output in a table so interested parties can examine them.

### Captured Whos??

Captured inputs and outputs are data points that are registered by either your system under test (stubs, mocks etc) or the output from your system its self.

Captured inputs and outputs are logged along side your test, for each test, so that interested parties can view them.

### That's all great, but still.. WHY would I want to log this stuff??

In BDD, a system's inputs and outputs are important to the business. Capturing the relationships between your data and the way the system handles can be demonstrated to your client. For example, your system could call a 3rd party, and you might want to model the interaction with stubs. 

GoGivens gives you a standardised way of rendering captured data alongside your tests so you don't have to worry about it.


### Rendered how?

The test framework parses your test file, and grabs the content. It strips all non-interesting parts out and leaves the Given/When/Then format in plain text ready for a GoGivensOutputGenerator to process the text. Interesting givens and Captured inputs and outputs are maps, which are rendered alongside your test givens as table data -- interesting givens are tablulated, and captured IO is listed.

A complete example of how to write a GoGivensOutputGenerator is given in the sister project [html spec](https://github.com/corbym/htmlspec) - written in Go.

## Example One - GoGivens in Practice <a name="example"></a>
```go
import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/testdata"
	"testing"
	"github.com/corbym/gocrest/is"
)

func TestMain(testmain *testing.M) {
	runOutput := testmain.Run()
	GenerateTestOutput() // You only need test main GenerateTestOutput() if you want to produce HTML output.
	os.Exit(runOutput)
}

func TestMyFirst(testing *testing.T) {
	Given(testing, someDataSetup).

		When(somethingHappens).

		Then(func(testing base.TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens) { // passed in testing should be used for assertions
		//do assertions
		then.AssertThat(testing, actual["actual"], is.EqualTo("some output"))
	})
}
```
Note you do not have to use "gocrest" assertions, you can still call all of testing.T's functions to fail the test or you can use any go testing assertion package compatible with testing.T.

When run, the above will produce an HTML output:

[Example Html](https://corbym.github.io/gogiven/example_test.shtml#github.com%2fcorbym%2fgogiven.TestMyFirst)

## Example Two - Table Tests <a name="tabletest-example"></a>

Table tests work the same way as normal go table tests. GoGivens will then mark which test failed, if they do, in your test output. 

Example:

```go
...
func TestMyFirst(testing *testing.T){
   var someRange = []struct {
		actual   string
		expected int
	}{
		{actual: "", expected: 0},
		{actual: "a", expected: 2},
	}
	for _, test := range someRange {
	   tst.Run(test.actual, func(weAreTesting *testing.T) {
	   	Given(weAreTesting, someDataSetup).
			When(someAction).
			Then(func(t TestingT, actual CapturedIO, givens InterestingGivens) {
			//do assertions
		AssertThat(t, actual.CapturedIO["actual"], is.EqualTo("some output"))
	   	})
	   }	
	}
}
...
```

The above test will still fail the test function as far as Go is concerned, but the test output will note that the iteration failed like this:

[Ranged Example Html](https://corbym.github.io/gogiven/example_test.shtml#github.com%2fcorbym%2fgogiven.TestMyFirst_Ranged)

**Note that comments are now rendered as "Noting that ..". In the above, the comment //do assertions would become "Noting that do assertions".**

### More Examples

* [Skipped test](https://corbym.github.io/gogiven/example_test.shtml#github.com%2fcorbym%2fgogiven.TestMyFirst_Skipped.func1)
* [Without a Given](https://corbym.github.io/gogiven/example_test.shtml#github.com%2fcorbym%2fgogiven.TestWithoutGiven)

# Content Generation <a name="content-gen"></a>

Gogivens comes defaultly configured with an html generator (```htmlspec.NewTestOutputGenerator```) that is consumed by a file generator (```generator.FileOutputGenerator```) (see the godoc for more information). The content generator implements the following interface:

```go
type GoGivensOutputGenerator interface {
	Generate(data *PageData) (output io.Reader)
	//ContentType is text/html, application/json or other mime type
	ContentType() string
}
```

The generated content ```(output io.Reader)``` is then consumed by an OutputListener:


```go
type OutputListener interface {
	Notify(testFilePath string, contentType string, output io.Reader)
}
```

If you want your own output listener just create your own and replace and/or append to the default listeners in your TestMain:

```go
func TestMain(testmain *testing.M) {
	gogiven.OutputListeners = []generator.OutputListener{new(MyFooListener)}
	// or alternately (or inclusively!)
	gogiven.OutputListeners = append(OutputListeners, new(MyBarListener))
	runOutput := testmain.Run()
	gogiven.GenerateTestOutput() // generates the output after the test is finished.
	os.Exit(runOutput)
}
```

## Setting the test file output (for the ```generator.FileOutputGenerator```)

You can add the environment variable GOGIVENS_OUTPUT_DIR to your env properties that points to a directory you want goGivens to report the test output to.

Default is the os's tmp directory.

## List of Pre-written Ouput Generators <a name="output-generator-list"></a>

GoGiven comes with the following output generators:

* HTML Spec: https://github.com/corbym/htmlspec - generates the output used in the test example. 
* JSON Spec: https://github.com/corbym/jsonspec - generates the output in JSON format. 
