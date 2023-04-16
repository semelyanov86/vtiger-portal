package service

import (
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/email/smtp"
)

type Services struct {
	Users UsersService
	Email smtp.Mailer
}

func NewServices(repos repository.Repositories, email smtp.Mailer) *Services {
	return &Services{
		Users: NewUsersService(repos.Users, repos.UsersCrm),
		Email: email,
	}
}
