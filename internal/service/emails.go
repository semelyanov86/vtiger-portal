package service

import (
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/email"
)

type EmailService struct {
	sender email.Sender
	config config.EmailConfig

	cache cache.Cache
}

type VerificationEmailInput struct {
	Name         string
	CompanyName  string
	SupportEmail string
	Email        string
}

type EmailServiceInterface interface {
	SendGreetingsToUser(input VerificationEmailInput) error
}

func NewEmailsService(sender email.Sender, config config.EmailConfig, cache cache.Cache) *EmailService {
	return &EmailService{
		sender: sender,
		config: config,
		cache:  cache,
	}
}

func (s EmailService) SendGreetingsToUser(input VerificationEmailInput) error {
	return s.sender.Send(input.Email, s.config.Templates.RegistrationEmail, input)
}

type MockEmailService struct {
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (s MockEmailService) SendGreetingsToUser(input VerificationEmailInput) error {
	return nil
}
