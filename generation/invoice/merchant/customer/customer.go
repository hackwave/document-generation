package customer

import (
	contact "../../contact"
	address "../../contact/address"
)

type Customer struct {
	FullName string
	Contacts contact.Contacts
	Address  address.Address
	Business string
}

func Individual(name string) *Customer {
	return &Customer{
		FullName: name,
	}
}

func Business(name, representative string) *Customer {
	return &Customer{
		Business: name,
		FullName: representative,
	}
}

func (self *Customer) Name() string {
	if len(self.Business) != 0 {
		return self.Business
	} else if len(self.FullName) != 0 {
		return self.FullName
	} else {
		return "n/a"
	}
}

func (self *Customer) IsBusiness() bool {
	return len(self.Business) != 0
}

func (self *Customer) BusinessName() string {
	if self.IsBusiness() {
		return self.Business
	} else {
		return "n/a"
	}
}

// TODO: Obvio move all the localbitcoins code together for dying out loud
func (self *Customer) AddLocalBitcoinsAccount(username string) *Customer {
	self.Contacts = append(self.Contacts, &contact.Contact{
		Type:    contact.Account,
		Service: "localbitcoins.com",
		Value:   username,
	})

	return self
}

func (self *Customer) LocalBitcoinsAccount() string {
	accountName := self.Contacts.AccountAt("localbitcoins.com")
	if 0 < len(accountName) {
		return accountName
	} else {
		return "n/a"
	}
}

func (self *Customer) GithubAccount() string {
	accountName := self.Contacts.AccountAt("localbitcoins.com")
	if 0 < len(accountName) {
		return accountName
	} else {
		return "n/a"
	}
}
