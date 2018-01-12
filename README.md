# gogiven
An alternative BDD spec framework for go. Builds on "go test" tool and builds on the go testing package.

[![Build status](https://travis-ci.org/corbym/gogiven.svg?branch=master)](https://github.com/corbym/gogiven)
[![Go Report Card](https://goreportcard.com/badge/github.com/corbym/gogiven)](https://goreportcard.com/report/github.com/corbym/gogiven)
[![GoDoc](https://godoc.org/github.com/corbym/gogiven?status.svg)](http://godoc.org/github.com/corbym/gogiven)

Inspired by [YATSPEC](https://github.com/bodar/yatspec).

## Example
```
import (
	. "github.com/corbym/gocrest/then"
	. "github.com/corbym/gogiven/given"
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
...
func someDataSetup(givens *InterestingGivens) {
    givens.Givens["1"] = "hi" //these keys and values will be displayed in the test output html in a table next to the test
    givens.Givens["2"] = "foo"
}
...
func someAction(capturedIo *CapturedIO, givens *InterestingGivens) {
    // call your functions here, feed output into capturedIO map - this will be displayed in the test output html
    capturedIo.CapturedIO["actual"] = new(T).callMyFun(givens.Givens["1"], givens.Givens["2"])
}
```
Note you do not have to use "gocrest" assertions, you can still call all of testing.T's functions to fail the test or you can use any go testing assertion package compatible with testing.T.

When run, the above will produce an HTML output:

[Example Html](http://htmlpreview.github.com/?https://raw.githubusercontent.com/corbym/gogiven/master/resources/example.html)

## Table Tests

Table tests work the same way as normal go table tests, but we use ThenFor which provides test context. goGivens will then mark which test in your loop failed. Example:

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
