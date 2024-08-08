package utils

var SupportedCurrencies = map[string]string{
	"USD": "USD",
	"EUR": "EUR",
	"CAD": "CAD",
}

func IsSupportedCurrencies(currency string) bool {
	_, ok := SupportedCurrencies[currency]
	return ok
}
