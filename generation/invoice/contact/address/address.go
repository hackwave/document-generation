package address

import (
	"fmt"
)

type Address struct {
	Street string 
	Zipcode string
	City string
	State string
	Country string
}

func (self Address) Format() string {
	return fmt.Sprintf("%s\n%s %s %s\n",
	  self.Street,
	  self.City,
	  self.State,
	  self.Zipcode,
	)
}
