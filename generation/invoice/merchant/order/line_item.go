package order

import (
	currency "../../currency"
	product "../product"
)

type LineItem struct {
	Product *product.Product

	ID string

	Name        string
	Description string

	Quantity float64
	Price    float64

	Currency currency.ISO
}

func (self *Order) AddLineItem(name string, quantity, price float64) *Order {
	self.LineItems = append(self.LineItems, &LineItem{
		Name:     name,
		Quantity: quantity,
		Price:    price,
		Currency: self.Currency,
	})
	return self
}

func (self *LineItem) Total() float64 {
	return self.Quantity * self.Price
}
