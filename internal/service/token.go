package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"time"
)

type TokensService struct {
	repo     repository.Tokens
	userRepo repository.Users
}

var ErrPasswordDoesNotMatch = errors.New("password does not match")

func NewTokensService(repo repository.Tokens, userRepo repository.Users) TokensService {
	return TokensService{
		repo:     repo,
		userRepo: userRepo,
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
