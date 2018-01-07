package gogiven

type CapturedIO struct {
	CapturedIO map[string]interface{}
}

func newCapturedIO() *CapturedIO {
	capturedIO := new(CapturedIO)
	capturedIO.CapturedIO = map[string]interface{}{}
	return capturedIO
}
