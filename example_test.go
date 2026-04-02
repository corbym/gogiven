package gogiven_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	. "github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven"
	"github.com/corbym/gogiven/base"
	"github.com/corbym/gogiven/testdata"
)

// TestGreetingService_PersonalisesGreeting verifies that the greeting service
// generates a personalised message for a registered user.
func TestGreetingService_PersonalisesGreeting(t *testing.T) {
	gogiven.Given(t, aRegisteredUserNamed("Alice")).
		When(aGreetingIsRequested).
		Then(func(t base.TestingT, captured testdata.CapturedIO, givens testdata.InterestingGivens) {
			// the greeting should address the user by their registered name
			AssertThat(t, captured["greeting"], is.EqualTo("Hello, Alice!"))
		})
}

func aRegisteredUserNamed(name string) func(givens testdata.InterestingGivens) {
	return func(givens testdata.InterestingGivens) {
		givens["userName"] = name
	}
}

func aGreetingIsRequested(captured testdata.CapturedIO, givens testdata.InterestingGivens) {
	name := givens["userName"].(string)
	captured["greeting"] = fmt.Sprintf("Hello, %s!", name)
}

// TestGreetingService_PersonalisesGreeting_ForManyUsers tests that the greeting service
// produces the correct personalised greeting for a range of user names.
// Comments on new lines will be split into paragraph sections.
//
// Each test case verifies that the greeting contains exactly the expected number of characters.
func TestGreetingService_PersonalisesGreeting_ForManyUsers(t *testing.T) {
	type greetingTestCase struct {
		userName       string
		expectedLength int
	}
	var testCases = []greetingTestCase{
		{userName: "Li", expectedLength: 10},   // "Hello, Li!"
		{userName: "Alice", expectedLength: 13}, // "Hello, Alice!"
	}
	for _, tc := range testCases {
		t.Run(tc.userName, func(tt *testing.T) {
			weAreTesting := base.NewTestMetaData(t.Name())
			gogiven.Given(weAreTesting, aRegisteredUserNamed(tc.userName)).
				When(aGreetingIsRequested).
				Then(func(t base.TestingT, captured testdata.CapturedIO, stored testdata.InterestingGivens) {
					// the greeting length should match the length of the formatted output
					AssertThat(t, captured["greeting"], has.Length(tc.expectedLength))
				})
		})
	}
}

// TestGreetingService_PersonalisesGreeting_SkipsUnknownLocale tests that the service
// skips producing a greeting for locales that are not yet supported.
func TestGreetingService_PersonalisesGreeting_SkipsUnknownLocale(t *testing.T) {
	type localeTestCase struct {
		userName string
		locale   string
	}
	var testCases = []localeTestCase{
		{userName: "Alice", locale: "en-US"},
		{userName: "Marie", locale: "fr-FR"},
	}
	for _, tc := range testCases {
		t.Run(tc.locale, func(t *testing.T) {
			gogiven.Given(t, aRegisteredUserNamed(tc.userName), withLocale(tc.locale)).
				SkippingThisOneIf(localeIsNotEnglish(tc.locale), "locale %s is not yet supported", tc.locale).
				When(aGreetingIsRequested).
				Then(func(t base.TestingT, captured testdata.CapturedIO, givens testdata.InterestingGivens) {
					AssertThat(t, captured["greeting"], is.EqualTo("Hello, Alice!"))
				})
		})
	}
}

func withLocale(locale string) func(givens testdata.InterestingGivens) {
	return func(givens testdata.InterestingGivens) {
		givens["locale"] = locale
	}
}

func localeIsNotEnglish(locale string) func(...interface{}) bool {
	return func(...interface{}) bool {
		return !strings.HasPrefix(locale, "en")
	}
}

// TestGreetingService_PersonalisesGreeting_NonDefaultParamName verifies that gogiven
// correctly parses test functions that use non-standard testing parameter names.
func TestGreetingService_PersonalisesGreeting_NonDefaultParamName(myT *testing.T) {
	gogiven.Given(myT, aRegisteredUserNamed("Alice")).
		When(aGreetingIsRequested).
		Then(func(thenT base.TestingT, captured testdata.CapturedIO, givens testdata.InterestingGivens) {
			// the greeting should address the user by their registered name
			AssertThat(thenT, captured["greeting"], is.EqualTo("Hello, Alice!"))
		})
}

// TestGreetingService_ReturnsDefaultGreeting verifies that the service returns a
// default greeting message when no user context is provided.
func TestGreetingService_ReturnsDefaultGreeting(t *testing.T) {
	gogiven.When(t, aDefaultGreetingIsRequested).
		Then(func(t base.TestingT, captured testdata.CapturedIO, givens testdata.InterestingGivens) {
			AssertThat(t, captured["greeting"], is.EqualTo("Hello, World!"))
		})
}

func aDefaultGreetingIsRequested(captured testdata.CapturedIO, givens testdata.InterestingGivens) {
	captured["greeting"] = "Hello, World!"
}
