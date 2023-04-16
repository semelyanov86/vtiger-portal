package service

import (
	"context"
	"errors"
	"github.com/octoper/go-ray"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type UserSignUpInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Code     string `json:"code" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=5,max=20"`
}

type UsersService struct {
	repo repository.Users
	crm  repository.UsersCrm
}

func NewUsersService(repo repository.Users, crm repository.UsersCrm) UsersService {
	return UsersService{repo: repo, crm: crm}
}

func (s UsersService) SignUp(ctx context.Context, input UserSignUpInput) (*domain.User, error) {
	users, err := s.crm.FindByEmail(ctx, input.Email)
	var user *domain.User
	if err != nil {
		return user, e.Wrap("can not find current email in crm", err)
	}

	for _, u := range users {
		if u.Code == input.Code {
			user = &u
		}
	}

	if user == nil || user.Crmid == "" {
		return user, e.Wrap("can not find user in vtiger", ErrUserNotFound)
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err = user.Password.Set(input.Password)
	if err != nil {
		return user, e.Wrap("can not hash password", err)
	}

	if err := s.repo.Insert(ctx, user); err != nil {
		return user, err
	}
	ray.Ray(user)
	return user, nil
}
