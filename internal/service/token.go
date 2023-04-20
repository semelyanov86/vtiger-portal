package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"time"
)

type TokensService struct {
	repo     repository.Tokens
	userRepo repository.Users
	emails   EmailServiceInterface
	config   config.Config
}

var ErrPasswordDoesNotMatch = errors.New("password does not match")

func NewTokensService(repo repository.Tokens, userRepo repository.Users, emails EmailServiceInterface, config config.Config) TokensService {
	return TokensService{
		repo:     repo,
		userRepo: userRepo,
		emails:   emails,
		config:   config,
	}
}

func (s TokensService) CreateAuthToken(ctx context.Context, login string, pass string) (*domain.Token, error) {
	user, err := s.userRepo.GetByEmail(ctx, login)
	if err != nil {
		return nil, e.Wrap("can not find user by email", err)
	}
	match := user.Password.Matches(pass)
	if !match {
		return nil, ErrPasswordDoesNotMatch
	}
	token, err := s.repo.New(ctx, user.Id, 24*time.Hour*90, domain.ScopeAuthentication)
	return token, err
}

func (s TokensService) SendPasswordResetToken(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if !user.IsActive {
		return ErrUserIsNotActive
	}
	token, err := s.repo.New(ctx, user.Id, 45*time.Minute, domain.ScopePasswordReset)
	if err != nil {
		return e.Wrap("Can not create new token", err)
	}
	emailData := PasswordRestoreData{
		Name:    user.FirstName + " " + user.LastName,
		Email:   user.Email,
		Token:   token.Plaintext,
		Valid:   token.Expiry,
		Company: s.config.Vtiger.Business.CompanyName,
		Support: s.config.Vtiger.Business.SupportEmail,
		Domain:  s.config.Domain,
		Subject: s.config.Email.Subjects.RestorePassword,
	}
	return s.emails.SendPasswordReset(emailData)
}
