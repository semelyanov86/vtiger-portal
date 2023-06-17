package repository

import (
	"context"
	"database/sql"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"time"
)

type NotificationsRepo struct {
	db *sql.DB
}

func NewNotificationsRepo(db *sql.DB) *NotificationsRepo {
	return &NotificationsRepo{
		db: db,
	}
}

func (r *NotificationsRepo) Insert(ctx context.Context, notification *domain.Notification) error {
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	var query = `
				INSERT INTO notifications (crmid, module, label, description, assigned_user_id, account_id, user_id, is_read, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	var args = []any{notification.Crmid, notification.Module, notification.Label, notification.Description, notification.AssignedUserId, notification.AccountId, notification.UserId, notification.IsRead, notification.CreatedAt, notification.UpdatedAt}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	notification.Id = id

	return nil
}

func (r *NotificationsRepo) GetNotificationsFromUserId(ctx context.Context, id string) ([]domain.Notification, error) {
	var query = `SELECT id, crmid, module, label, description, assigned_user_id, account_id, user_id, is_read, created_at, updated_at FROM notifications WHERE user_id = ? AND is_read = 0`
	var notifications = make([]domain.Notification, 0)
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification domain.Notification
		err = rows.Scan(&notification.Id, &notification.Crmid, &notification.Module, &notification.Label, &notification.Description, &notification.AssignedUserId, &notification.AccountId, &notification.UserId, &notification.IsRead, &notification.CreatedAt, &notification.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationsRepo) MarkNotificationAsRead(ctx context.Context, id int64, userId string) error {
	var query = `UPDATE notifications SET is_read = 1, updated_at = NOW() WHERE id = ? AND user_id = ?`
	var args = []any{id, userId}
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
