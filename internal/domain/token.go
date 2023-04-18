package domain

import "time"

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
