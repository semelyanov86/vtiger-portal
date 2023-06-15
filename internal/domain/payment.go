package domain

import "time"

type PaymentConfig struct {
	PublishableKey string `json:"publishableKey"`
}

type Payment struct {
	ID              int64     `json:"id"`
	StripePaymentId string    `json:"stripe_payment_id"`
	UserId          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	PaymentMethod   string    `json:"payment_method"`
	Status          int       `json:"status"`
	ParentId        string    `json:"parent_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
