package gogiven

import (
	"testing"
	"runtime"
	"sync"
)

type CapturedIO struct {
	CapturedIO map[string]interface{}
}

type InterestingGivens struct {
	Givens map[string]interface{}
}

var someTests *sync.Map

func Given(testing *testing.T, given ...func(givens *InterestingGivens)) *some {
	if someTests == nil {
		someTests = new(sync.Map)
	}

	funcProgramCounters := make([]uintptr, 1)
	noOfEntries := runtime.Callers(2, funcProgramCounters)
	if noOfEntries == 0 {
		panic("eek")
	}
	// get the info of the actual function that's in the pointer
	function := runtime.FuncForPC(funcProgramCounters[0] - 1)
	if function == nil {
		panic("arrgh")
	}

	keyFor := keyFor(someTests, function.Name())
	some := newSome(newTestMetaData(testing, keyFor), function, funcProgramCounters, given...)
	someTests.Store(keyFor, some)
	return some
}

func keyFor(somes *sync.Map, name string) string {
	if _, ok := somes.Load(name); !ok {
		return name
	}
	return keyFor(somes, name+"_1")
}
