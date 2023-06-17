package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type NotificationsCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewNotificationsCrm(config config.Config, cache cache.Cache) NotificationsCrm {
	return NotificationsCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (n NotificationsCrm) GetAllNotifications(ctx context.Context) ([]domain.Notification, error) {
	// TODO: Implement receiving notification from CRM
	return []domain.Notification{}, nil
}

func (n NotificationsCrm) DeleteReceivedNotification(ctx context.Context, crmid string) error {
	// TODO: Implement clearing notification in crm
	return nil
}
