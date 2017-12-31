package given

import (
	"testing"
	"runtime"
	"fmt"
	"io/ioutil"
)

type CapturedIO struct {
	CapturedIO map[string]interface{}
}

type InterestingGivens struct {
	Givens map[string]interface{}
}

type Some struct {
	testingT          *testing.T
	InterestingGivens *InterestingGivens
	CapturedIO        *CapturedIO
}

func Given(testing *testing.T, given ...func(givens *InterestingGivens) *InterestingGivens) *Some {
	some := new(Some)
	some.testingT = testing

	some.CapturedIO = new(CapturedIO)
	some.CapturedIO.CapturedIO = map[string]interface{}{}

	givens := new(InterestingGivens)
	givens.Givens = map[string]interface{}{}

	if len(given) > 0 {
		for _, someGiven := range given {
			some.InterestingGivens = someGiven(givens)
		}
	} else {
		some.InterestingGivens = givens
	}
	return some
}

func (some *Some) When(action func(actual *CapturedIO, givens *InterestingGivens) *CapturedIO) *Some {
	some.CapturedIO = action(some.CapturedIO, some.InterestingGivens)

	fpcs := make([]uintptr, 1)
	n := runtime.Callers(2, fpcs)
	if n == 0 {
		panic("eek")
	}
	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		panic("arrgh")
	}

	generateHtml(fun, fpcs, some.testingT)
	return some
}

func (some *Some) Then(assertions func(actual *CapturedIO, givens *InterestingGivens)) *Some {
	assertions(some.CapturedIO, some.InterestingGivens)
	return some
}

func (some *Some) ElseFail() {
	if some.CapturedIO == nil {
		some.testingT.FailNow()
	}
}

func generateHtml(runtimeCaller *runtime.Func, fpcs []uintptr, testing *testing.T) {
	file, line := runtimeCaller.FileLine(fpcs[0] - 1)
	_, err := ioutil.ReadFile(file)
	if err != nil {
		panic("file not found")
	}
	//fmt.Print(string(dat))

	fmt.Printf("%s called in: %s at line %d", runtimeCaller.Name(), file, line)
}
