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
func someDataSetup(givens *InterestingGivens) *InterestingGivens {
    givens.Givens["1"] = "hi"
    givens.Givens["2"] = "foo"
    return givens
}
...
func someAction(capturedIo *CapturedIO, givens *InterestingGivens) *CapturedIO {
    // call your functions here, feed output into capturedIO map
    capturedIo.CapturedIO["actual"] = new(T).callMyFun(givens.Givens["1"], givens.Givens["2"])
    return capturedIo
}
```

When run, the above will produce an HTML output:

(TBD)