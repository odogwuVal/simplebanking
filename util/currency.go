package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	NGA = "NGA"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, NGA:
		return true
	}
	return false
}
