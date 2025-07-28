package domain

type CurrencyPair struct {
	From   string
	To     string
	Amount float64
}

type Quote struct {
	Provider string
	Rate     float64
	Err      error
}
