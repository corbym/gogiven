package generator

import (
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven/base"
)

func TestNewPageData(t *testing.T) {
	someMap := &base.SomeMap{}
	(*someMap)["foo"] = base.NewSome(t, "foob", new(base.TestMetaData), []string{"givenwhenthen"})
	then.AssertThat(t, NewPageData("foo", someMap).TestResults, is.Not(is.Nil()))
}
