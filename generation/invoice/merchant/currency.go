package merchant

var UYUtoUSDfor2019 = 35.2838353413654

// NOTE: For 2019 Rates
func ConvertUYUtoUSD(uyuAmount float64) float64 {
	return (uyuAmount / UYUtoUSDfor2019)
}
