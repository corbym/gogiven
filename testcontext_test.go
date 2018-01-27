package gogiven_test

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven"
	"testing"
)

const someFilename = "foofar.go"

func TestNewTestContext(t *testing.T) {
	context := gogiven.NewTestContext(someFilename)
	then.AssertThat(t, context.FileName(), is.EqualTo(someFilename))
	then.AssertThat(t, context.SomeTests(), has.TypeName("*gogiven.safeMap"))
}
