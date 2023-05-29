package domain

import (
	"github.com/semelyanov86/vtiger-portal/internal/utils"
	"github.com/semelyanov86/vtiger-portal/pkg/validator"
	"time"
)
import "github.com/jameskeane/bcrypt"

type User struct {
	Id                 int64     `json:"id"`
	Crmid              string    `json:"crmid"`
	FirstName          string    `json:"firstname"`
	LastName           string    `json:"lastname"`
	Description        string    `json:"description"`
	AccountId          string    `json:"account_id"`
	AccountName        string    `json:"account_name"`
	Title              string    `json:"title"`
	Department         string    `json:"department"`
	Email              string    `json:"email"`
	Password           Password  `json:"-"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	IsActive           bool      `json:"is_active"`
	MailingCity        string    `json:"mailingcity"`
	MailingStreet      string    `json:"mailingstreet"`
	MailingCountry     string    `json:"mailingcountry"`
	OtherCountry       string    `json:"othercountry"`
	MailingState       string    `json:"mailingstate"`
	MailingPoBox       string    `json:"mailingpobox"`
	OtherCity          string    `json:"othercity"`
	OtherState         string    `json:"otherstate"`
	MailingZip         string    `json:"mailingzip"`
	OtherZip           string    `json:"otherzip"`
	OtherStreet        string    `json:"otherstreet"`
	OtherPoBox         string    `json:"otherpobox"`
	Image              string    `json:"image"`
	Version            int       `json:"-"`
	Code               string    `json:"-"`
	Imageattachmentids string    `json:"imageattachmentids"`
	Imagecontent       string    `json:"imagecontent"`
	Phone              string    `json:"phone"`
	AssignedUserId     string    `json:"assigned_user_id"`
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
	v.Check(user.Crmid != "", "data.crmid", "must be provided")
	ValidateEmail(v, user.Email)
	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}
	if user.Password.Hash == nil {
		panic("missing password Hash for user")
	}
}

func ConvertMapToUser(m map[string]any) User {
	user := User{}

	for k, v := range m {
		switch k {
		case "id":
			user.Crmid = v.(string)
		case "firstname":
			user.FirstName = v.(string)
		case "lastname":
			user.LastName = v.(string)
		case "description":
			user.Description = v.(string)
		case "account_id":
			user.AccountId = v.(string)
		case "account_name":
			user.AccountName = v.(string)
		case "title":
			user.Title = v.(string)
		case "department":
			user.Department = v.(string)
		case "email":
			user.Email = v.(string)
		case "created_at":
			user.CreatedAt, _ = time.Parse(time.RFC3339, v.(string))
		case "updated_at":
			user.UpdatedAt, _ = time.Parse(time.RFC3339, v.(string))
		case "is_active":
			user.IsActive = v.(bool)
		case "mailingcity":
			user.MailingCity = v.(string)
		case "mailingstreet":
			user.MailingStreet = v.(string)
		case "mailingcountry":
			user.MailingCountry = v.(string)
		case "othercountry":
			user.OtherCountry = v.(string)
		case "mailingstate":
			user.MailingState = v.(string)
		case "mailingpobox":
			user.MailingPoBox = v.(string)
		case "othercity":
			user.OtherCity = v.(string)
		case "otherstate":
			user.OtherState = v.(string)
		case "mailingzip":
			user.MailingZip = v.(string)
		case "otherzip":
			user.OtherZip = v.(string)
		case "otherstreet":
			user.OtherStreet = v.(string)
		case "otherpobox":
			user.OtherPoBox = v.(string)
		case "image":
			user.Image = v.(string)
		case "imageattachmentids":
			user.Imageattachmentids = v.(string)
		case "phone":
			user.Phone = v.(string)
		case "assigned_user_id":
			user.AssignedUserId = v.(string)
		}
	}

	return user
}

func (u User) ConvertToMap() (map[string]any, error) {
	result, err := utils.ConvertStructToMap(u)
	if err != nil {
		return result, err
	}

	return result, nil
}
