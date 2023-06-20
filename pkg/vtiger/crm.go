package vtiger

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"strconv"
	"strings"
)

type PaginationQueryFilter struct {
	Page     int
	PageSize int
	Client   string
	Contact  string
	Parent   string
	Sort     string
	Filters  map[string]any
	Search   string
}

type QueryFieldsProps struct {
	DefaultSort  string
	SearchFields []string
	ClientField  string
	AccountField string
	TableName    string
}

func (c VtigerConnector) GetAll(ctx context.Context, filter PaginationQueryFilter, fields QueryFieldsProps) ([]map[string]any, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM " + fields.TableName + " WHERE "
	sort := filter.Sort
	if sort == "" {
		sort = fields.DefaultSort
	}
	if filter.Search != "" {
		for _, field := range fields.SearchFields {
			query += " " + field + " LIKE '%" + filter.Search + "%' OR"
		}
		query = strings.Trim(query, "OR")
	} else {
		if fields.AccountField != "" {
			query += fields.AccountField + " = " + filter.Client + " "
		}
		if fields.ClientField != "" {
			query += "OR " + fields.ClientField + " = " + filter.Client + " "
		}
	}
	query += GenerateOrderByClause(sort)
	query += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"
	result, err := c.Query(ctx, query)
	if err != nil {
		return nil, e.Wrap("can not execute query "+query+", got error", err)
	}
	return result.Result, nil
}

func (c VtigerConnector) GetByWhereClause(ctx context.Context, filter PaginationQueryFilter, field string, value string, table string) ([]map[string]any, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM " + table + " WHERE " + field + " = '" + value + "' LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"
	result, err := c.Query(ctx, query)
	if err != nil {
		return nil, e.Wrap("can not execute query "+query+", got error: ", err)
	}
	return result.Result, nil
}
