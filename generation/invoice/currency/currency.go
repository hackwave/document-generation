package currency

type ISO int

const (
	USD ISO = iota
	UYU
	BTC
)

func (self ISO) String() string {
	switch self {
	case USD:
		return "USD"
	case UYU:
		return "UYU"
	case BTC:
		return "BTC"
	default:
		return "USD"
	}
}

func MarshalISO(iso string) ISO {
	switch iso {
	case USD.String():
		return USD
	case UYU.String():
		return UYU
	case BTC.String():
		return BTC
	default:
		return USD
	}
}
