package service

import (
	"context"
	"encoding/json"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/logger"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
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
	vtiger     vtiger.VtigerConnector
}

const (
	PENDING = iota
	SUCCEEDED
	CANCELLED
	PROCESSING
	REQUIRES_ACTION
	REQUIRES_CAPTURE
	REQUIRES_CONFIRMATION
	REQUIRES_PAYMENT_METHOD
	CREATED
)

type PaymentIntent struct {
	Currency          string  `json:"currency"`
	PaymentMethodType string  `json:"paymentMethodType"`
	SoId              string  `json:"so_id"`
	InvoiceId         string  `json:"invoice_id"`
	Amount            float64 `json:"amount"`
	UserId            string  `json:"userId"`
	AccountId         string  `json:"accountId"`
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
		vtiger:     vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
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
		AccountId:       req.AccountId,
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

func (p Payments) AcceptPayment(ctx context.Context, e stripe.Event) (domain.Payment, error) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(e.Data.Raw, &paymentIntent)
	if err != nil {
		return domain.Payment{}, err
	}
	payment, err := p.repository.GetByStripeId(ctx, paymentIntent.ID)
	if err != nil {
		return payment, err
	}
	payment.Status = p.getNumericIntentStatus(paymentIntent)

	return p.repository.UpdatePayment(ctx, payment)
}

func (p Payments) ConfirmPayment(ctx context.Context, event stripe.PaymentIntent) (domain.Payment, error) {
	payment, err := p.repository.GetByStripeId(ctx, event.ID)
	if err != nil {
		return payment, e.Wrap("can not get payment by stripe id "+event.ID, err)
	}
	payment.Status = p.getNumericIntentStatus(event)

	if payment.Status == SUCCEEDED {
		go func() {
			err := p.reviseModuleSuccessStatus(payment.ParentId)
			if err != nil {
				logger.Error(logger.GenerateErrorMessageFromString(err.Error()))
			}
		}()
	}
	return p.repository.UpdatePayment(ctx, payment)
}

func (p Payments) reviseModuleSuccessStatus(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m := make(map[string]any)
	m["id"] = id
	m["invoicestatus"] = p.config.Payment.PaidInvoiceStatus
	m["sostatus"] = p.config.Payment.PaidSoStatus
	_, err := p.vtiger.Revise(ctx, m)
	return err
}

func (p Payments) GetPayments(ctx context.Context, id string) ([]domain.Payment, error) {
	return p.repository.GetPaymentsFromAccountId(ctx, id)
}

func (p Payments) getNumericIntentStatus(intent stripe.PaymentIntent) int {
	switch intent.Status {
	case stripe.PaymentIntentStatusCanceled:
		return CANCELLED
	case stripe.PaymentIntentStatusProcessing:
		return PROCESSING
	case stripe.PaymentIntentStatusRequiresAction:
		return REQUIRES_ACTION
	case stripe.PaymentIntentStatusRequiresCapture:
		return REQUIRES_CAPTURE
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		return REQUIRES_PAYMENT_METHOD
	case stripe.PaymentIntentStatusSucceeded:
		return SUCCEEDED
	}
	return CREATED
}
