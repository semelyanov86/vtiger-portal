package domain

import (
	"github.com/semelyanov86/vtiger-portal/pkg/validator"
	"time"
)

const ScopeActivation = "activation"
const ScopeAuthentication = "authentication"
const ScopePasswordReset = "password-reset"

type Token struct {
	ID        int64     `json:"id"`
	Plaintext string    `json:"token"`
	Hash      string    `json:"-"`
	UserId    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
