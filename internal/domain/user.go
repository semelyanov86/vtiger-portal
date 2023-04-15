package domain

import (
	"github.com/semelyanov86/vtiger-portal/pkg/validator"
	"time"
)
import "github.com/jameskeane/bcrypt"

type User struct {
	Id             int64     `json:"id"`
	Crmid          int64     `json:"crmid"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Description    string    `json:"description"`
	AccountId      int       `json:"account_id"`
	AccountName    string    `json:"account_name"`
	Title          string    `json:"title"`
	Department     string    `json:"department"`
	Email          string    `json:"email"`
	Password       Password  `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsActive       bool      `json:"is_active"`
	MailingCity    string    `json:"mailingcity"`
	MailingStreet  string    `json:"mailingstreet"`
	MailingCountry string    `json:"mailingcountry"`
	OtherCountry   string    `json:"othercountry"`
	MailingState   string    `json:"mailingstate"`
	MailingPoBox   string    `json:"mailingpobox"`
	OtherCity      string    `json:"othercity"`
	OtherState     string    `json:"otherstate"`
	MailingZip     string    `json:"mailingzip"`
	OtherZip       string    `json:"otherzip"`
	OtherStreet    string    `json:"otherstreet"`
	OtherPoBox     string    `json:"otherpobox"`
	Image          string    `json:"image"`
	Version        int       `json:"-"`
}

var AnonymousUser = &User{}

type Password struct {
	Plaintext *string
	Hash      []byte
}

func (p *Password) Set(plaintext string) error {
	hash, err := bcrypt.HashBytes([]byte(plaintext))
	if err != nil {
		return err
	}

	p.Plaintext = &plaintext
	p.Hash = hash

	return nil
}

func (p *Password) Matches(plaintext string) bool {
	return bcrypt.MatchBytes([]byte(plaintext), p.Hash)
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "data.attributes.email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "data.attributes.email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "data.attributes.password", "must be provided")
	v.Check(len(password) > 8, "data.attributes.password", "must be at least 8 bytes long")
	v.Check(len(password) < 72, "data.attributes.password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Email != "", "data.email", "must be provided")
	v.Check(user.Crmid != 0, "data.crmid", "must be provided")
	ValidateEmail(v, user.Email)
	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}
	if user.Password.Hash == nil {
		panic("missing password Hash for user")
	}
}
