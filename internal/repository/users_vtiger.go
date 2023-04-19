package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type UsersVtiger struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewUsersVtiger(config config.Config, cache cache.Cache) UsersVtiger {
	return UsersVtiger{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (receiver UsersVtiger) FindByEmail(ctx context.Context, email string) ([]domain.User, error) {
	result, err := receiver.vtiger.Query(ctx, "SELECT * FROM Contacts WHERE email = '"+email+"';")

	users := make([]domain.User, 0)
	if err != nil {
		return users, e.Wrap("can not get user by email "+email, err)
	}

	for _, m := range result.Result {
		curUser := domain.ConvertMapToUser(m)
		curUser.Code = m[receiver.config.Vtiger.Business.CodeField].(string)
		users = append(users, curUser)
	}
	return users, nil
}

func (receiver UsersVtiger) RetrieveById(ctx context.Context, id string) (domain.User, error) {
	result, err := receiver.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.User{}, e.Wrap("can not retrieve user with id"+id, err)
	}
	user := domain.ConvertMapToUser(result.Result)
	return user, nil
}
