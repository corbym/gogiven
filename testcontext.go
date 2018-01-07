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

func newTestMetaData(t *testing.T, testName string) *TestingT {
	testContext := new(TestingT)
	testContext.t = t
	testContext.TestName = testName
	return testContext
}

func (ctx *TestingT) Skipped() {
	t := ctx.t
	t.Skipped()
}

func (ctx *TestingT) Parallel() {
	t := ctx.t
	t.Helper()
	t.Parallel()
}

func (ctx *TestingT) Logf(format string, args ...interface{}) {
	t := ctx.t
	t.Helper()
	ctx.TestOutput = fmt.Sprintf(format, args...)
	ctx.failed = true
	t.Logf(format, args...)
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
	ctx.t.Helper()
}
