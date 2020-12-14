package order

import (
	"fmt"
	"time"

	currency "../../currency"
	customer "../customer"
)

var orderDateFormat = "2006-01-02"

func ParseDate(date string) time.Time {
	parsedDate, _ := time.Parse("2006-01-02", date)
	return parsedDate
}

type Order struct {
	ID       int
	IDPrefix string

	Merchant interface{}
	Customer *customer.Customer

	Timestamp time.Time
	Currency  currency.ISO

	Data map[string]map[string]string

	LineItems []*LineItem
}

func (self *Order) Total() (total float64) {
	for _, item := range self.LineItems {
		total += item.Total()
	}
	return total
}

func (self *Order) CreatedAt() string {
	return fmt.Sprintf(self.Timestamp.Format(orderDateFormat))
}

func (self *Order) InvoiceID() string {
	return fmt.Sprintf("%s-00%v", self.IDPrefix, self.ID)
}
