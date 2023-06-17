package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
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

func (receiver UsersVtiger) FindContactsInAccount(ctx context.Context, filter PaginationQueryFilter) ([]string, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT id FROM Contacts WHERE account_id = " + filter.Client + " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"
	result, err := receiver.vtiger.Query(ctx, query)
	ids := make([]string, 0)
	if err != nil {
		return nil, e.Wrap("can not execute query "+query+", got error", err)
	}
	for _, data := range result.Result {
		ids = append(ids, data["id"].(string))
	}
	return ids, nil
}

func (receiver UsersVtiger) ClearUserCodeField(ctx context.Context, id string) (domain.User, error) {
	result, err := receiver.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.User{}, e.Wrap("can not retrieve user with id"+id, err)
	}
	data := result.Result
	codeField := receiver.config.Vtiger.Business.CodeField
	if codeField != "" {
		data[codeField] = ""
	}
	result, err = receiver.vtiger.Update(ctx, data)
	if err != nil {
		return domain.User{}, e.Wrap("can not retrieve user with id"+id, err)
	}
	user := domain.ConvertMapToUser(result.Result)
	return user, nil
}

func (receiver UsersVtiger) Update(ctx context.Context, id string, user domain.User) (domain.User, error) {
	userMap, err := user.ConvertToMap()
	if err != nil {
		return user, e.Wrap("can not convert to map", err)
	}
	userMap["id"] = id

	result, err := receiver.vtiger.Revise(ctx, userMap)
	if err != nil {
		return user, e.Wrap("can send update map to vtiger", err)
	}
	user = domain.ConvertMapToUser(result.Result)
	return user, nil
}

func (receiver UsersVtiger) RetrieveContactMap(ctx context.Context, id string) (map[string]any, error) {
	result, err := receiver.vtiger.Retrieve(ctx, id)
	if err != nil {
		return map[string]any{}, e.Wrap("can not retrieve user with id"+id, err)
	}
	return result.Result, nil
}

func (receiver UsersVtiger) ChangeSettingField(ctx context.Context, id string, field string, value bool) error {
	input := make(map[string]any)
	input["id"] = id
	if value {
		input[field] = "1"
	} else {
		input[field] = "0"
	}
	_, err := receiver.vtiger.Revise(ctx, input)
	return err
}
