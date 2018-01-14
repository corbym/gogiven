package gogiven_test

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven"
	"testing"
)

const fileThatDoesNotExist = "foofar.go"

func TestNewTestContext_Panics(t *testing.T) {
	defer func() {
		rcv := recover()
		then.AssertThat(t, rcv, is.Not(is.Nil()))
	}()
	gogiven.NewTestContext(fileThatDoesNotExist)
}
