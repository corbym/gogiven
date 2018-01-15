package gogiven_test

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/corbym/gogiven"
	"github.com/corbym/gogiven/testdata"
	"os"
	"testing"
)

func TestGenerateTestOutput(t *testing.T) {
	defer func() {
		remove := os.Remove(ofFileInTmpDir("generator_test.html"))
		then.AssertThat(t, remove, is.Nil())
	}()
	t.Parallel()
	gogiven.Given(t, func(givens testdata.InterestingGivens) {

	})
	gogiven.GenerateTestOutput()
	then.AssertThat(t, fileExists(ofFileInTmpDir("generator_test.html")), inTmpDir())
}

func TestGenerateTestOutput_DefaultsToCurrentDir(t *testing.T) {
	defer func() {
		remove := os.Remove("./generator_test.html")
		then.AssertThat(t, remove, is.Nil())
	}()
	defer func() { os.Setenv("GOGIVENS_OUTPUT_DIR", "") }()
	os.Setenv("GOGIVENS_OUTPUT_DIR", "doesnotexist")
	gogiven.Given(t)
	gogiven.GenerateTestOutput()
	then.AssertThat(t, fileExists("./generator_test.html"), inTmpDir())
}
