package email

import (
	"github.com/semelyanov86/vtiger-portal/pkg/validator"
)

const (
	minEmailLen = 3
	maxEmailLen = 255
)

var emailRegex = validator.EmailRX

func IsEmailValid(email string) bool {
	if len(email) < minEmailLen || len(email) > maxEmailLen {
		return false
	}

	return emailRegex.MatchString(email)
}
