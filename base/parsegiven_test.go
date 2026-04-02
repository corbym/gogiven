package base_test

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven/base"
	"testing"
)

const exampleTestFileName = "../example_test.go"

func TestParseGivenWhenThen_Panics(t *testing.T) {
	defer func() {
		rcv := recover()
		AssertThat(t, rcv, is.Not(is.Nil()))
	}()
	base.ParseGivenWhenThen("foo", "Arfg")
}

func TestParseGivenWhenThen_TextOutputContent(testing *testing.T) {
	g := base.ParseGivenWhenThen(".TestMyFirst", exampleTestFileName)

	AssertThat(testing, g.GivenWhenThen, has.Length(5))

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("Given the system setup"))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("When something happens"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("Then"))
	AssertThat(testing, g.GivenWhenThen[3], is.EqualTo("Noting that we do some assertions here commenting why"))
	AssertThat(testing, g.GivenWhenThen[4], is.EqualTo("Assert that captured \"actual\" is equal to \"some output\""))
}

func TestParseGivenWhenThen_WithoutGiven(testing *testing.T) {
	g := base.ParseGivenWhenThen(".TestWithoutGiven", exampleTestFileName)

	AssertThat(testing, g.GivenWhenThen, has.Length(3))

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("When something happens"))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("Then"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("Assert that actual \"actual\" is equal to \"some output\""))
}

func TestParseGivenWhenThen_PanicsWithoutGivenOrWhen(testing *testing.T) {
	defer func() {
		recover := recover()
		AssertThat(testing, recover, is.Not(is.Nil()))
	}()
	base.ParseGivenWhenThen(".TestParseGivenWhenThen_PanicsWithoutGivenOrWhen", "./given_test.go")
}

func TestParseGivenWhenThen_FuncWithReturnType(testing *testing.T) {
	g := base.ParseGivenWhenThen(".TestMyFirst_Skipped", exampleTestFileName)
	AssertThat(testing, g.GivenWhenThen, has.Length(5))

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("Given the system setup that is a bit dodgy to test"))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("Skipping this one if the value is fff test \"some data % s does not work yet\" test actual"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("When something happens with the test"))
	AssertThat(testing, g.GivenWhenThen[3], is.EqualTo("Then"))
	AssertThat(testing, g.GivenWhenThen[4], is.EqualTo("Assert that test actual is equal to \"a\" reason \"we only want to assert if test actual is a not empty\""))
}

func TestParseGivenWhenThen_NonDefaultTParamName(testing *testing.T) {
	g := base.ParseGivenWhenThen(".TestMyFirst_NonDefaultParamName", exampleTestFileName)

	AssertThat(testing, g.GivenWhenThen, has.Length(5))

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("Given the system setup"))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("When something happens"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("Then"))
	AssertThat(testing, g.GivenWhenThen[3], is.EqualTo("Noting that we do some assertions here commenting why"))
	AssertThat(testing, g.GivenWhenThen[4], is.EqualTo("Assert that captured \"actual\" is equal to \"some output\""))
}

func TestParseGivenWhenThen_RangedTextOutput(testing *testing.T) {
	parsedTest := base.ParseGivenWhenThen(".TestMyFirst_Ranged", exampleTestFileName)
	AssertThat(testing, parsedTest.GivenWhenThen, has.Length(5))

	AssertThat(testing, parsedTest.GivenWhenThen[0], is.EqualTo("Given the system setup with test data test"))
	AssertThat(testing, parsedTest.GivenWhenThen[1], is.EqualTo("When something happens with the test"))
	AssertThat(testing, parsedTest.GivenWhenThen[2], is.EqualTo("Then"))
	AssertThat(testing, parsedTest.GivenWhenThen[3], is.EqualTo("Noting that do assertions"))
	AssertThat(testing, parsedTest.GivenWhenThen[4], is.EqualTo("Assert that stored \"actual\" has length test expected"))
}

func TestParseGivenWhenThen_IncludesComment(testing *testing.T) {
	givenWhenThen := base.ParseGivenWhenThen(".TestMyFirst_Ranged", exampleTestFileName)

	AssertThat(testing, givenWhenThen.Comment[0], is.EqualTo("This test tests over a range of values. Lorum ipsum dolor, lorum ipsum dolor lorum ipsum dolor. Lorum ipsum dolor."))
	AssertThat(testing, givenWhenThen.Comment[1], is.EqualTo("Comments on new lines will be split into paragraph sections."))
}
