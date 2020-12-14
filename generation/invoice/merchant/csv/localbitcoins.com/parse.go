package localbitcoins

import (
	"strconv"
	"strings"
	"time"

	currency "../../../currency"
	customer "../../customer"
	order "../../order"
)

var orderDateFormat = "2006-01-02"

func OrderFromCSVRecord(record []string) *order.Order {
	// record[0] => id
	// record[1] => created_at
	customerLocalBitcoinsUsername := record[2] // => buyer (userbintc
	// record[3] => seller
	orderTradeType := record[4] // => trade_type
	// record[5] => btc_amount
	// record[6] => btc_traded
	// record[7] => fee_btc
	// record[8] => btc_amount_less_fee (NOTE: NICE FUCKING NAMING)
	orderBitcoinAmount := record[9] // => btc_final
	//orderCurrencyAmount := record[10] // => fiat_amount
	// record[11] => fiat_fee
	orderBitcoinPrice := record[12] // => fiat_per_btc
	//orderCurrencyPerBitcoin := record[12]
	orderCurrencyName := record[13] // => currency_name
	// record[14] => exchange_rate
	//orderCurrency := record[14] // =>
	// record[15] => transaction_released_at
	orderCreatedAt := record[15]
	// record[16] => online_provider
	// record[17] => reference
	orderLocalBitcoinsReference := record[17]
	// NOTE: These had to be manually appended because the data
	// LocalBitcoins.com gives is really not well designed.
	// record[18] => customer
	customerFullName := record[18]
	// record[19] => comapny
	customerBusinessName := record[19]

	if len(orderTradeType) == 11 && orderTradeType[8:] == "SELL" {
		return nil
	}

	createdAt, _ := time.Parse(orderDateFormat, strings.Split(orderCreatedAt, " ")[0])

	bitcoinAmount, _ := strconv.ParseFloat(orderBitcoinAmount, 4)
	bitcoinPrice, _ := strconv.ParseFloat(orderBitcoinPrice, 4)
	// Total
	//currencyAmount, _ := strconv.ParseFloat(orderCurrencyAmount, 4)

	data := map[string]map[string]string{
		"localbitcoins.com": map[string]string{
			"reference": orderLocalBitcoinsReference,
			"username":  customerLocalBitcoinsUsername,
		},
	}

	// TODO: Eventually store all customers inside of merchant. Then Create
	// orders FROM the customer only.
	var c *customer.Customer
	if len(customerBusinessName) != 0 {
		c = customer.Individual(customerFullName)
	} else {
		c = customer.Business(customerBusinessName, customerFullName)
	}

	currencyISO := currency.MarshalISO(orderCurrencyName)

	return &order.Order{
		ID:        0000,
		IDPrefix:  "BTCDONAT",
		Timestamp: createdAt,
		Currency:  currencyISO,
		Customer:  c,
		Data:      data,
		LineItems: []*order.LineItem{
			&order.LineItem{
				ID:          "DONATION-01",
				Name:        "Bitcoin",
				Description: "Donation obtained for contributing to open source projects",
				Quantity:    bitcoinAmount,
				Price:       bitcoinPrice,
				Currency:    currencyISO,
			},
		},
	}
}
