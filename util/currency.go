package util

const (
	EURO = "EURO"
	USD  = "USD"
	CAD  = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case EURO, USD, CAD:
		return true
	}
	return false
}
