package generator

import (
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven/base"
)

func TestNewPageData(t *testing.T) {
	someMap := base.SomeMap{}
	someMap["foo"] = base.NewSome(t, "foob",
		&base.TestMetaData{TestID: "narf"},
		base.ParsedTestContent{GivenWhenThen: []string{"givenwhenthen"}},
	)
	then.AssertThat(t, NewPageData("foo", &someMap).TestData, is.Not(is.Nil[[]TestData]()))
}
