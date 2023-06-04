package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"time"
)

var MockedUser = domain.User{
	Id:          1,
	Crmid:       "12x11",
	AccountId:   "11x1",
	FirstName:   "Sergey",
	LastName:    "Emelyanov",
	Description: "Test Description",
	Email:       "emelyanov86@km.ru",
	Password:    domain.Password{},
	CreatedAt:   time.Time{},
	UpdatedAt:   time.Time{},
	IsActive:    true,
	Version:     1,
}

type UsersMock struct {
}

func NewUsersMock() *UsersMock {
	return &UsersMock{}
}

func (r *UsersMock) Insert(ctx context.Context, user *domain.User) error {
	user.Id = 1
	user.Version = 1
	return nil
}

func (r *UsersMock) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return MockedUser, ErrRecordNotFound
}

func (r *UsersMock) Update(ctx context.Context, user *domain.User) error {
	return nil
}

func (r *UsersMock) GetById(ctx context.Context, id int64) (domain.User, error) {
	return MockedUser, nil
}

func (r *UsersMock) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	return &MockedUser, nil
}

func (r *UsersMock) SaveOtp(ctx context.Context, otpSecret string, otpUrl string, userId int64) error {
	return nil
}

func (r *UsersMock) EnableAndVerifyOtp(ctx context.Context, userId int64) error {
	return nil
}

func (r *UsersMock) VerifyOrInvalidateOtp(ctx context.Context, userId int64, valid bool) error {
	return nil
}

func (r *UsersMock) DisableOtp(ctx context.Context, userId int64) error {
	return nil
}
