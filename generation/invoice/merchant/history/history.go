package history

import (
	order "../order"
)

type History struct {
	OrderID int

	Orders []*order.Order
}

func (self *History) IterateOrderID() int {
	self.OrderID += 1
	return self.OrderID
}

// Order History //////////////////////////////////////////////////////////////
func (self *History) AddOrder(o *order.Order) *History {
	self.Orders = append(self.Orders, o)
	return self
}

func (self *History) HasOrder(o *order.Order) bool {
	for _, order := range self.Orders {
		if order.ID == o.ID {
			return true
		}
	}
	return false
}

func (self *History) AtIndex(index int) *order.Order {
	for _, order := range self.Orders {
		if order.ID == index {
			return order
		}
	}
	return nil
}

func (self *History) WithID(id string) *order.Order {
	for _, order := range self.Orders {
		if order.InvoiceID() == id {
			return order
		}
	}
	return nil
}
