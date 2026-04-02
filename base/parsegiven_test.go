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
	g := base.ParseGivenWhenThen(".TestGreetingService_PersonalisesGreeting", exampleTestFileName)

	AssertThat(testing, g.GivenWhenThen, has.Length(5))

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("Given a registered user named \"alice\""))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("When a greeting is requested"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("Then"))
	AssertThat(testing, g.GivenWhenThen[3], is.EqualTo("Noting that the greeting should address the user by their registered name"))
	AssertThat(testing, g.GivenWhenThen[4], is.EqualTo("Assert that captured \"greeting\" is equal to \"hello alice !\""))
}

func TestParseGivenWhenThen_WithoutGiven(testing *testing.T) {
	g := base.ParseGivenWhenThen(".TestGreetingService_ReturnsDefaultGreeting", exampleTestFileName)

	AssertThat(testing, g.GivenWhenThen, has.Length(3))

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("When a default greeting is requested"))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("Then"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("Assert that captured \"greeting\" is equal to \"hello world !\""))
}

func TestParseGivenWhenThen_PanicsWithoutGivenOrWhen(testing *testing.T) {
	defer func() {
		recover := recover()
		AssertThat(testing, recover, is.Not(is.Nil()))
	}()
	base.ParseGivenWhenThen(".TestParseGivenWhenThen_PanicsWithoutGivenOrWhen", "./given_test.go")
}

func TestParseGivenWhenThen_FuncWithReturnType(testing *testing.T) {
	g := base.ParseGivenWhenThen(".TestGreetingService_PersonalisesGreeting_SkipsUnknownLocale", exampleTestFileName)
	AssertThat(testing, g.GivenWhenThen, has.Length(5))

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("Given a registered user named tc user name with locale tc locale"))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("Skipping this one if locale is not english tc locale \"locale % s is not yet supported\" tc locale"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("When a greeting is requested"))
	AssertThat(testing, g.GivenWhenThen[3], is.EqualTo("Then"))
	AssertThat(testing, g.GivenWhenThen[4], is.EqualTo("Assert that captured \"greeting\" is equal to \"hello alice !\""))
}

func TestParseGivenWhenThen_NonDefaultTParamName(testing *testing.T) {
	g := base.ParseGivenWhenThen(".TestGreetingService_PersonalisesGreeting_NonDefaultParamName", exampleTestFileName)

	AssertThat(testing, g.GivenWhenThen, has.Length(5))

	AssertThat(testing, g.GivenWhenThen[0], is.EqualTo("Given a registered user named \"alice\""))
	AssertThat(testing, g.GivenWhenThen[1], is.EqualTo("When a greeting is requested"))
	AssertThat(testing, g.GivenWhenThen[2], is.EqualTo("Then"))
	AssertThat(testing, g.GivenWhenThen[3], is.EqualTo("Noting that the greeting should address the user by their registered name"))
	AssertThat(testing, g.GivenWhenThen[4], is.EqualTo("Assert that captured \"greeting\" is equal to \"hello alice !\""))
}

func TestParseGivenWhenThen_RangedTextOutput(testing *testing.T) {
	parsedTest := base.ParseGivenWhenThen(".TestGreetingService_PersonalisesGreeting_ForManyUsers", exampleTestFileName)
	AssertThat(testing, parsedTest.GivenWhenThen, has.Length(5))

	AssertThat(testing, parsedTest.GivenWhenThen[0], is.EqualTo("Given a registered user named tc user name"))
	AssertThat(testing, parsedTest.GivenWhenThen[1], is.EqualTo("When a greeting is requested"))
	AssertThat(testing, parsedTest.GivenWhenThen[2], is.EqualTo("Then"))
	AssertThat(testing, parsedTest.GivenWhenThen[3], is.EqualTo("Noting that the greeting length should match the length of the formatted output"))
	AssertThat(testing, parsedTest.GivenWhenThen[4], is.EqualTo("Assert that captured \"greeting\" has length tc expected length"))
}

func TestParseGivenWhenThen_IncludesComment(testing *testing.T) {
	givenWhenThen := base.ParseGivenWhenThen(".TestGreetingService_PersonalisesGreeting_ForManyUsers", exampleTestFileName)

	AssertThat(testing, givenWhenThen.Comment[0], is.EqualTo("TestGreetingService_PersonalisesGreeting_ForManyUsers tests that the greeting service"))
	AssertThat(testing, givenWhenThen.Comment[1], is.EqualTo("produces the correct personalised greeting for a range of user names."))
	AssertThat(testing, givenWhenThen.Comment[2], is.EqualTo("Comments on new lines will be split into paragraph sections."))
}
