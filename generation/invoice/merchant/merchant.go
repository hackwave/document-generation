package merchant

import (
	"fmt"
	"time"

	contact "../contact"
	address "../contact/address"
	currency "../currency"
	csv "./csv"
	localbitcoins "./csv/localbitcoins.com"
	customer "./customer"
	history "./history"
	order "./order"
)

type Merchant struct {
	Name   string
	Slogan string

	Contacts contact.Contacts

	PhysicalAddress address.Address
	WebsiteAddress  string

	Currency currency.ISO
	Country  string

	History *history.History

	Data map[string]string
}

func (self *Merchant) NewOrder(c *customer.Customer) *order.Order {
	order := &order.Order{
		ID:        self.History.IterateOrderID(),
		Merchant:  self,
		Customer:  c,
		Timestamp: time.Now(), // Default
		Currency:  self.Currency,
	}

	self.History.AddOrder(order)

	return order
}

func (self *Merchant) OrderByID(id string) *order.Order {
	fmt.Println("id string:", id)
	return self.History.WithID(id)
}

func (self *Merchant) OrderByIndex(index int) *order.Order {
	if len(self.History.Orders) >= index {
		return self.History.Orders[index]
	} else {
		return nil
	}
}

// TODO:
// OrderFromLocalBitcoinsRecord(record []string)
func Default(name string) *Merchant {
	switch name {
	case "hackwave laboratories":
		hackwave := &Merchant{
			Name:     "Hackwave Laboratories",
			Slogan:   "Open Source Hardware & Software Engineering",
			Currency: currency.USD,
			Country:  "USA",
			Contacts: contact.Contacts{
				&contact.Contact{Type: contact.Email, Value: "contact@hackwave.org"},
				&contact.Contact{Type: contact.Account, Service: "Github.com", Value: "hackwave"},
			},
			PhysicalAddress: address.Address{
				Street:  "123 Fake Street.",
				City:    "Springfield",
				State:   "OR",
				Zipcode: "12345",
			},
			WebsiteAddress: "hackwave.org",
			Data:           map[string]string{},
			History:        &history.History{},
		}

		var temporaryOrder *order.Order
		for _, record := range csv.LoadFile("./data/hackwave/sales.csv") {
			temporaryOrder = localbitcoins.OrderFromCSVRecord(record)
			newOrder := &order.Order{
				ID:        hackwave.History.IterateOrderID(),
				IDPrefix:  temporaryOrder.IDPrefix,
				Timestamp: temporaryOrder.Timestamp,
				Currency:  temporaryOrder.Currency,
				Customer:  temporaryOrder.Customer,
				Data:      temporaryOrder.Data,
				LineItems: temporaryOrder.LineItems,
			}

			hackwave.History.Orders = append(hackwave.History.Orders, newOrder)

		}

		return hackwave
	default:
		return nil
	}
}

func (self *Merchant) OrdersTotal() (total float64) {
	for _, order := range self.History.Orders {
		if order.Currency == currency.UYU {
			total += ConvertUYUtoUSD(order.Total())
		} else {
			total += order.Total()
		}
	}
	return total
}
