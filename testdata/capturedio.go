package testdata

//CapturedIO contains the captured inputs and outputs for the test.
// These could be the interactions between the system and the stubbed endpoints,
// or any value that the system under test produces.
// Underlying the CapturedIO is a map[string]interface{}.
// CapturedIO objects are provided to the test by the When() and Then() functions.
type CapturedIO struct {
	CapturedIO map[string]interface{}
}

func NewCapturedIO() *CapturedIO {
	capturedIO := new(CapturedIO)
	capturedIO.CapturedIO = map[string]interface{}{}
	return capturedIO
}
