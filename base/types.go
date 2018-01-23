package base

import (
	"github.com/corbym/gogiven/testdata"
)

// GivenData is a func type that gets given Some interesting givens as a parameter
type GivenData func(givens testdata.InterestingGivens)

// CapturedIOData is a func type that gets given a reference to Some CapturedIO data for the test
type CapturedIOData func(capturedIO testdata.CapturedIO)

//CapturedIOGivenData is a combination of GivenData and CapturedIOData types
type CapturedIOGivenData func(capturedIO testdata.CapturedIO, givens testdata.InterestingGivens)

// TestingWithGiven gives a func declaration including testingT, CapturedIO and InterestingGivens parameters.
type TestingWithGiven func(testingT TestingT, actual testdata.CapturedIO, givens testdata.InterestingGivens)
