package gogiven
import (
	. "github.com/corbym/gocrest/then"
	"github.com/corbym/gocrest/is"
	"testing"
)

func TestMyFirst(testing *testing.T){
	Given(testing, someDataSetup).
		When(someAction).
		Then(func(testing TestingT, actual *CapturedIO, givens *InterestingGivens) {
		//do assertions
		AssertThat(testing, actual.CapturedIO["actual"], is.EqualTo("some output"))
	})
}

func someDataSetup(givens *InterestingGivens) {
	givens.Givens["1"] = "hi" //these keys and values will be displayed in the test output html in a table next to the test
	givens.Givens["2"] = "foo"
}
func someAction(capturedIo *CapturedIO, givens *InterestingGivens) {
	// call your functions here, feed output into capturedIO map - this will be displayed in the test output html
	capturedIo.CapturedIO["actual"] = callMyFun(givens.Givens["1"], givens.Givens["2"])
}

func callMyFun(some1 interface{}, some2 interface{}) interface{} {
	return "some output"
}
