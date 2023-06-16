package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"time"
)

type PaymentsRepo struct {
	db *sql.DB
}

func NewPaymentsRepo(db *sql.DB) *PaymentsRepo {
	return &PaymentsRepo{
		db: db,
	}
}

func (r *PaymentsRepo) Insert(ctx context.Context, payment *domain.Payment) error {
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	var query = `
				INSERT INTO payments (stripe_payment_id, user_id, amount, currency, payment_method, status, created_at, updated_at, parent_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	var args = []any{payment.StripePaymentId, payment.UserId, payment.Amount, payment.Currency, payment.PaymentMethod, payment.Status, payment.CreatedAt, payment.UpdatedAt, payment.ParentId}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	payment.ID = id

	return nil
}

func (r *PaymentsRepo) GetByStripeId(ctx context.Context, id string) (domain.Payment, error) {
	var query = `SELECT id, stripe_payment_id, user_id, amount, currency, payment_method, status, created_at, updated_at, parent_id FROM payments WHERE stripe_payment_id = ?`
	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, query, id).Scan(&payment.ID, &payment.StripePaymentId, &payment.UserId, &payment.Amount, &payment.Currency, &payment.PaymentMethod, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt, &payment.ParentId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return payment, ErrRecordNotFound
		default:
			return payment, err
		}
	}
	return payment, nil
}

func (r *PaymentsRepo) UpdatePayment(ctx context.Context, payment domain.Payment) (domain.Payment, error) {
	var query = `UPDATE payments SET stripe_payment_id = ?, user_id = ?, amount = ?, currency = ?, payment_method = ?, status = ?, updated_at = NOW(), parent_id = ?`
	var args = []any{payment.StripePaymentId, payment.UserId, payment.Amount, payment.Currency, payment.PaymentMethod, payment.Status, payment.ParentId}
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return payment, err
	}
	payment.UpdatedAt = time.Now()
	return payment, nil
}
