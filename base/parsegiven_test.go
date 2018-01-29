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

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("Given testing the system setup"))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("When something happens"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("Then noting that passed in testing should be used for assertions"))
	AssertThat(testing, g.GivenWhenThen[3], is.EqualTo("Noting that we do some assertions here commenting why"))
	AssertThat(testing, g.GivenWhenThen[4], is.EqualTo("Assert that testing the data returned \"actual\" is equal to \"some output\""))
}

func TestParseGivenWhenThen_WithoutGiven(testing *testing.T) {
	g := base.ParseGivenWhenThen(".TestWithoutGiven", exampleTestFileName)

	AssertThat(testing, g.GivenWhenThen, has.Length(3))

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("When t something happens"))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("Then"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("Assert that testing actual \"actual\" is equal to \"some output\""))
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

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("Given we are testing the system setup that is a bit dodgy to test"))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("Skipping this one if the value is fff test \"some data % s does not work yet\" test actual"))
	// TODO: fix this
	//AssertThat(testing, givenWhenThen[4], is.EqualTo("some data % s does not work yet \"test actual\""))
}

func TestParseGivenWhenThen_RangedTextOutput(testing *testing.T) {
	parsedTest := base.ParseGivenWhenThen(".TestMyFirst_Ranged", exampleTestFileName)
	AssertThat(testing, parsedTest.GivenWhenThen, has.Length(5))

	AssertThat(testing, parsedTest.GivenWhenThen[0], is.EqualTo("Given we are testing the system setup with test data test"))
	AssertThat(testing, parsedTest.GivenWhenThen[1], is.EqualTo("When something happens with the test"))
	AssertThat(testing, parsedTest.GivenWhenThen[2], is.EqualTo("Then"))
	AssertThat(testing, parsedTest.GivenWhenThen[3], is.EqualTo("Noting that do assertions"))
}

func TestParseGivenWhenThen_IncludesComment(testing *testing.T) {
	givenWhenThen := base.ParseGivenWhenThen(".TestMyFirst_Ranged", exampleTestFileName)

	AssertThat(testing, givenWhenThen.Comment[0], is.EqualTo("This test tests over a range of values"))
	AssertThat(testing, givenWhenThen.Comment[1], is.EqualTo("Do not remove this comment"))
}
