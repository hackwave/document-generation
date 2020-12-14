package product

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
}

func Define(id, name, description string) *Product {
	return &Product{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

func (self *Product) SetPrice(price float64) *Product {
	self.Price = price
	return self
}
