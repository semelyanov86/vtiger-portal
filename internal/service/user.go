package service

import (
	"context"
	"errors"
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
	_, err := s.repo.GetByEmail(ctx, input.Email)
	if !errors.Is(repository.ErrRecordNotFound, err) {
		return nil, repository.ErrDuplicateEmail
	}
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
	FillVtigerContactWithAdditionalValues(user, input.Password)
	if user == nil || user.Crmid == "" {
		return user, e.Wrap("can not find user in vtiger", ErrUserNotFound)
	}

	if err := s.repo.Insert(ctx, user); err != nil {
		return user, err
	}

	return user, nil
}

func FillVtigerContactWithAdditionalValues(user *domain.User, password string) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := user.Password.Set(password)
	if err != nil {
		return e.Wrap("can not hash password", err)
	}
	return nil
}
