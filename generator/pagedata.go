package generator

import "github.com/corbym/gogiven/base"

//PageData is the struct that populates the template with data from the test output.
type PageData struct {
	Title   string
	SomeMap *base.SomeMap
}
