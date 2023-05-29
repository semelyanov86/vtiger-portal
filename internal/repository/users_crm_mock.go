package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
)

type UsersCrmMock struct {
	user domain.User
}

func NewUsersCrmMock(user domain.User) UsersCrmMock {
	return UsersCrmMock{
		user: user,
	}
}

func (receiver UsersCrmMock) FindByEmail(ctx context.Context, email string) ([]domain.User, error) {
	return []domain.User{
		receiver.user,
	}, nil
}

func (receiver UsersCrmMock) RetrieveById(ctx context.Context, id string) (domain.User, error) {
	return receiver.user, nil
}

func (receiver UsersCrmMock) ClearUserCodeField(ctx context.Context, id string) (domain.User, error) {
	return receiver.user, nil
}

func (receiver UsersCrmMock) FindContactsInAccount(ctx context.Context, filter PaginationQueryFilter) ([]string, error) {
	return []string{"17x16"}, nil
}

func (receiver UsersCrmMock) Update(ctx context.Context, id string, user domain.User) (domain.User, error) {
	return receiver.user, nil
}
