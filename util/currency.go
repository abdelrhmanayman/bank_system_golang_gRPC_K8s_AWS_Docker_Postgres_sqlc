package util

const (
	USD = "USD"
	EGP = "EGP"
	EUR = "EUR"
	SEK = "SEK"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, EGP, EUR, SEK:
		return true
	default:
		return false
	}
}
