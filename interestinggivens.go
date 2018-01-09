package gogiven

//InterestingGivens contains a map which can hold the values used to set up the test.
// You can add actual vs. expected output for example.
// InterestingGivens is created by the test framework and passed into the When() and Then()
// functions.
type InterestingGivens struct {
	Givens map[string]interface{} ``
}

func newInterestingGivens() *InterestingGivens {
	givens := new(InterestingGivens)
	givens.Givens = map[string]interface{}{}
	return givens
}
