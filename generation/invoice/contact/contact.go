package contact

import (
	"strings"
)

type Contacts []*Contact

type Type int

const (
	Phone Type = iota
	IM
	Email
	Account
)

type Contact struct {
	Type    Type
	Service string
	Value   string
}

func (self Contacts) ContactType(contactType Type) (typedContacts Contacts) {
	for _, contact := range self {
		if contact.Type == contactType {
			typedContacts = append(typedContacts, contact)
		}
	}
	return typedContacts
}

func (self Contacts) AccountAt(serviceName string) string {
	for _, contact := range self.ContactType(Account) {
		if strings.ToLower(contact.Service) == serviceName {
			return contact.Value
		}
	}
	return ""
}

func (self Contacts) Email() string {
	emails := self.ContactType(Email)
	if len(emails) > 0 {
		email := emails[0]
		if email != nil {
			return email.Value
		}
	}
	return ""
}

func (self Contacts) Phone() string {
	phones := self.ContactType(Phone)
	if len(phones) > 0 {
		phone := phones[0]
		if phone != nil {
			return phone.Value
		}
	}
	return ""
}
