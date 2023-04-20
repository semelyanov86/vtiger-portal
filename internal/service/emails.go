package service

import (
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/email"
	"time"
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

type PasswordRestoreData struct {
	Name    string
	Token   string
	Valid   time.Time
	Company string
	Support string
	Domain  string
	Email   string
	Subject string
}

type EmailServiceInterface interface {
	SendGreetingsToUser(input VerificationEmailInput) error
	SendPasswordReset(input PasswordRestoreData) error
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

func (s EmailService) SendPasswordReset(input PasswordRestoreData) error {
	return s.sender.Send(input.Email, s.config.Templates.RestorePasswordEmail, input)
}

type MockEmailService struct {
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (s MockEmailService) SendGreetingsToUser(input VerificationEmailInput) error {
	return nil
}

func (s MockEmailService) SendPasswordReset(input PasswordRestoreData) error {
	return nil
}
