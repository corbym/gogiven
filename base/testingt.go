package base

//TestingT is a convenience interface that matches some methods of `testing.T`
type TestingT interface {
	Failed() bool
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
	Name() string
}
