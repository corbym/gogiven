package model

import "github.com/corbym/gogiven/generator"

// JSONData is the JSON model for the test content.
type JSONData struct {
	Title     string       `json:"title"`
	TestState []*testState `json:"test_state"`
}

// NewJSONData creates a new json data object for marshalling test data.
func NewJSONData(pageData generator.PageData) *JSONData {
	jsonPageData := new(JSONData)
	jsonPageData.Title = pageData.Title
	for _, v := range pageData.TestData {
		jsonPageData.TestState = append(jsonPageData.TestState, newTestState(v))
	}
	return jsonPageData
}
