package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"time"
)

type Notifications struct {
	cache          cache.Cache
	config         config.Config
	manager        ManagerService
	repository     repository.NotificationsRepo
	crm            repository.NotificationsCrm
	userRepository repository.Users
}

var cachedAccounts map[string][]domain.User

func NewNotificationsService(cache cache.Cache, config config.Config, manager ManagerService, repository repository.NotificationsRepo, crm repository.NotificationsCrm, userRepository repository.Users) Notifications {
	return Notifications{
		cache:          cache,
		config:         config,
		manager:        manager,
		repository:     repository,
		crm:            crm,
		userRepository: userRepository,
	}
}

func (n Notifications) ImportNotifications() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	notifications, err := n.crm.GetAllNotifications(ctx)
	if err != nil {
		return e.Wrap("can not get notifications from crm", err)
	}
	for _, notification := range notifications {
		if notification.AccountId != "" {
			contacts, ok := cachedAccounts[notification.AccountId]
			if !ok {
				contacts, err = n.userRepository.GetAllByAccountId(ctx, notification.AccountId)
				if err != nil {
					return e.Wrap("can not get related contacts for account "+notification.AccountId, err)
				}
				cachedAccounts[notification.AccountId] = contacts
			}
			for _, contact := range contacts {
				notification.UserId = contact.Crmid
				err = n.repository.Insert(ctx, &notification)
				if err != nil {
					return e.Wrap("can not insert in database with id "+notification.Crmid+" for user"+contact.Crmid, err)
				}
			}

			err = n.crm.DeleteReceivedNotification(ctx, notification.Crmid)
			if err != nil {
				return e.Wrap("can not insert delete notification in crm "+notification.Crmid, err)
			}
		}
	}
	return nil
}

func (n Notifications) GetNotificationsByUserId(ctx context.Context, userId string) ([]domain.Notification, error) {
	notifications, err := n.repository.GetNotificationsFromUserId(ctx, userId)
	if err != nil {
		return notifications, e.Wrap("can not get notifications", err)
	}
	result := make([]domain.Notification, 0)
	for _, notification := range notifications {
		manager, err := n.manager.GetManagerById(ctx, notification.AssignedUserId)
		if err != nil {
			return notifications, e.Wrap("can not get manager "+notification.AssignedUserId, err)
		}
		notification.Manager = manager
		result = append(result, notification)
	}
	return result, nil
}

func (n Notifications) MarkNotificationRead(ctx context.Context, id int64, userId string) error {
	return n.repository.MarkNotificationAsRead(ctx, id, userId)
}
