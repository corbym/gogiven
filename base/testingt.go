package base

// TestingT is a convenience interface that matches Some methods of `testing.T`
type TestingT interface {
	Failed() bool
	Fail()
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
	Name() string
	Skipf(format string, args ...interface{})
}
