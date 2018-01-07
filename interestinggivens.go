package gogiven

type InterestingGivens struct {
	Givens map[string]interface{} ``
}

func newInterestingGivens() *InterestingGivens {
	givens := new(InterestingGivens)
	givens.Givens = map[string]interface{}{}
	return givens
}
