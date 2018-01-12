# gogiven
An alternative BDD spec framework for go. Builds on "go test" tool and builds on the go testing package.

[![Build status](https://travis-ci.org/corbym/gogiven.svg?branch=master)](https://github.com/corbym/gogiven)
[![Go Report Card](https://goreportcard.com/badge/github.com/corbym/gogiven)](https://goreportcard.com/report/github.com/corbym/gogiven)
[![GoDoc](https://godoc.org/github.com/corbym/gogiven?status.svg)](http://godoc.org/github.com/corbym/gogiven)

Inspired by [YATSPEC](https://github.com/bodar/yatspec).

## Table of Contents
1. [Introduction](#introduction)
2. [Example One - GoGivens in practice](#example)
3. [Example Two - Table Tests](#tabletest-example)
4. [Setting the test file output](<a name="file-output-settings"></a>)

## Introduction <a name="introduction"></a>

Go Givens is a lightweight TDD framework for producing test specifications directly from the code you write.

Go Givens parses your test file and produces a human-readable output in a specified directory, containing all the tests, captured data and other related information regarding your test such as success or failure.

### Interesting whats?

Interesting givens are data points that you want to use in your tests that are important for it to function correctly. "Given" data is then used by your application to fulfil some sort of request or process.

Interesting givens are generated as output along side your test output in a table so interested parties can examine them.

### Captured Whos??

Captured inputs and outputs are data points that are registered by either your system under test (stubs, mocks etc) or the output from your system its self.

Captured inputs and outputs are logged along side your test, for each test, so that interested parties can view them.

## Example One 0 GoGivens in Practice <a name="example"></a>
```
import (
	. "github.com/corbym/gocrest/then"
	. "github.com/corbym/gogiven"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest"
)

func TestMain(testmain *testing.M) {
	runOutput := testmain.Run()
	GenerateTestOutput() // You only need test main GenerateTestOutput() if you want to produce HTML output.
	os.Exit(runOutput)
}

func TestMyFirst(testing *testing.T){
   Given(testing, someDataSetup).
        When(someAction).
        Then(func(testing TestingT, actual *CapturedIO, givens *InterestingGivens) { // passed in testing should be used for assertions
        //do assertions
        AssertThat(testing, actual.CapturedIO["actual"], is.EqualTo("some output"))
    })
}
```
Note you do not have to use "gocrest" assertions, you can still call all of testing.T's functions to fail the test or you can use any go testing assertion package compatible with testing.T.

When run, the above will produce an HTML output:

[Example Html](http://htmlpreview.github.com/?https://raw.githubusercontent.com/corbym/gogiven/master/resources/example.html)

## Example Two - Table Tests <a name="tabletest-example"></a>

Table tests work the same way as normal go table tests. GoGivens will then mark which test in your loop failed. Example:

```
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
	   Given(testing, someDataSetup).
		When(someAction).
		Then(func(t *TestingT, actual *CapturedIO, givens *InterestingGivens) {
		//do assertions
		AssertThat(t, actual.CapturedIO["actual"], is.EqualTo("some output"))
	    })
	}
}
...
```
This will still fail the test function as far as Go is concerned, but the test output will note that the iteration failed like this:

[Ranged Example Html](http://htmlpreview.github.com/?https://raw.githubusercontent.com/corbym/gogiven/master/resources/example2.html)


## Setting the test file output <a name="file-output-settings"></a>

You can add the environment variable GOGIVENS_OUTPUT_DIR to your env properties that points to a direectory you want goGivens to report the test output to.

Default is the os's tmp directory.