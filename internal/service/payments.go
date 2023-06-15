package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	_ "github.com/stripe/stripe-go/v72/webhook"
	"runtime"
	"time"
)

type Payments struct {
	cache      cache.Cache
	config     config.Config
	currency   CurrencyService
	repository repository.PaymentsRepo
}

const (
	PENDING = iota
	SUCCEEDED
	FAILED
)

type PaymentIntent struct {
	Currency          string  `json:"currency"`
	PaymentMethodType string  `json:"paymentMethodType"`
	SoId              string  `json:"so_id"`
	InvoiceId         string  `json:"invoice_id"`
	Amount            float64 `json:"amount"`
	UserId            string  `json:"userId"`
}

func NewPaymentsService(cache cache.Cache, config config.Config, currency CurrencyService, repository repository.PaymentsRepo) Payments {
	stripe.Key = config.Payment.StripeKey

	// For sample support and debugging, not required for production:
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "customer-portal",
		Version: runtime.Version(),
		URL:     config.Domain,
	})

	return Payments{
		cache:      cache,
		config:     config,
		currency:   currency,
		repository: repository,
	}
}

func (p Payments) CreatePaymentIntent(ctx context.Context, req PaymentIntent) (*stripe.PaymentIntent, error) {
	var formattedPaymentMethodType []*string

	if req.PaymentMethodType == "link" {
		formattedPaymentMethodType = append(formattedPaymentMethodType, stripe.String("link"), stripe.String("card"))
	} else {
		formattedPaymentMethodType = append(formattedPaymentMethodType, stripe.String(req.PaymentMethodType))
	}

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(req.Amount * 100)),
		Currency:           stripe.String(req.Currency),
		PaymentMethodTypes: formattedPaymentMethodType,
	}

	// If this is for an ACSS payment, we add payment_method_options to create
	// the Mandate.
	if req.PaymentMethodType == "acss_debit" {
		params.PaymentMethodOptions = &stripe.PaymentIntentPaymentMethodOptionsParams{
			ACSSDebit: &stripe.PaymentIntentPaymentMethodOptionsACSSDebitParams{
				MandateOptions: &stripe.PaymentIntentPaymentMethodOptionsACSSDebitMandateOptionsParams{
					PaymentSchedule: stripe.String("sporadic"),
					TransactionType: stripe.String("personal"),
				},
			},
		}
	}

	res, err := paymentintent.New(params)
	if err != nil {
		return res, err
	}
	var parentId = req.SoId
	if parentId == "" {
		parentId = req.InvoiceId
	}
	payment := domain.Payment{
		StripePaymentId: res.ID,
		UserId:          req.UserId,
		Amount:          req.Amount,
		Currency:        req.Currency,
		PaymentMethod:   req.PaymentMethodType,
		Status:          PENDING,
		ParentId:        parentId,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = p.repository.Insert(ctx, &payment)
	return res, err
}

func (p Payments) AcceptPayment(ctx context.Context, e stripe.Event) error {
	payment, err := p.repository.GetByStripeId(ctx, e.ID)
	if err != nil {
		return err
	}
	payment.Status = SUCCEEDED

	return p.repository.UpdatePayment(ctx, payment)
}
