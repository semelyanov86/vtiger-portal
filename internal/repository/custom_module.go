package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/utils"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
	"strings"
)

type CustomModuleCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewCustomModuleCrm(config config.Config, cache cache.Cache) CustomModuleCrm {
	return CustomModuleCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (m CustomModuleCrm) GetAll(ctx context.Context, filter PaginationQueryFilter, custom vtiger.Module, fields []string) ([]map[string]any, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM " + custom.Name + " WHERE "
	sort := filter.Sort
	accountField := m.findAccountField(custom)
	contactField := m.findContactField(custom)
	if filter.Search != "" {
		for i, field := range fields {
			if i == 0 {
				continue
			}
			query += field + " LIKE '%" + filter.Search + "%' OR"
		}
		query = strings.Trim(query, "OR")
		query += " "
	} else {
		if accountField != nil {
			query += accountField.Name + " = " + filter.Client + " OR "
		}
		if contactField != nil {
			query += contactField.Name + " = " + filter.Client + " "
		}
		query = strings.Trim(query, "OR")
		query += " "
	}
	query += GenerateOrderByClause(sort)
	query += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"

	records := make([]map[string]any, 0)

	result, err := m.vtiger.Query(ctx, query)
	if err != nil {
		return nil, e.Wrap("can not execute query "+query+", got error", err)
	}
	for _, data := range result.Result {
		if accountField != nil && data[accountField.Name] == filter.Client {
			records = append(records, data)
			continue
		}
		if contactField != nil && data[contactField.Name] == filter.Client {
			records = append(records, data)
		}
	}
	return records, nil
}

func (m CustomModuleCrm) Count(ctx context.Context, client string, contact string, custom vtiger.Module) (int, error) {
	accountField := m.findAccountField(custom)
	contactField := m.findContactField(custom)
	body := make(map[string]string)
	if accountField != nil && client != "" {
		body[accountField.Name] = client
	}
	if contactField != nil && contact != "" {
		body[contactField.Name] = client
	}
	return m.vtiger.Count(ctx, custom.Name, body)
}

func (m CustomModuleCrm) findAccountField(module vtiger.Module) *vtiger.ModuleField {
	return utils.FindFieldByRefers(module, "Accounts")
}

func (m CustomModuleCrm) findContactField(module vtiger.Module) *vtiger.ModuleField {
	return utils.FindFieldByRefers(module, "Contacts")
}
