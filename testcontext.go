package gogiven

import (
	"fmt"
	"testing"
)

type TestingT struct {
	TestName   string
	failed     bool
	TestOutput string
	t          *testing.T
}

func (ctx *TestingT) Logf(format string, args ...interface{}) {
	t := ctx.t
	t.Helper()
	ctx.TestOutput = fmt.Sprintf(format, args...)
	ctx.failed = true
	ctx.t.Logf(format, args...)
}

func (ctx *TestingT) Errorf(format string, args ...interface{}) {
	t := ctx.t
	t.Helper()
	ctx.TestOutput = fmt.Sprintf(format, args...)
	ctx.failed = true
	t.Errorf(format, args...)
}

func (ctx *TestingT) FailNow() {
	t := ctx.t
	t.Helper()
	ctx.failed = true
	t.FailNow()
}
func (ctx *TestingT) HasFailed() bool {
	return ctx.failed
}

func (ctx *TestingT) Helper() {
	i := ctx.t.Helper
	i()
}
