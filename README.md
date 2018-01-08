# gogiven
An alternative BDD spec framework for go.

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

func TestMyFirst(testing *testing.T){
   Given(testing, someDataSetup).
        When(someAction).
        Then(func(actual *CapturedIO, givens *InterestingGivens) {
        //do assertions
        AssertThat(testing, actual.CapturedIO["actual"], is.EqualTo("some output"))
    })
}
...
func someDataSetup(givens *InterestingGivens) {
    givens.Givens["1"] = "hi"
    givens.Givens["2"] = "foo"
}
...
func someAction(capturedIo *CapturedIO, givens *InterestingGivens) {
    // call your functions here, feed output into capturedIO map
    capturedIo.CapturedIO["actual"] = new(T).callMyFun(givens.Givens["1"], givens.Givens["2"])
}
```

When run, the above will produce an HTML output:

(TBD)

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

(TBD)
