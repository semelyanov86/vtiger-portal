package vtiger

import (
	"context"
)

type Connector interface {
	GetAll(ctx context.Context, filter PaginationQueryFilter, fields QueryFieldsProps) ([]map[string]any, error)
	GetByWhereClause(ctx context.Context, filter PaginationQueryFilter, field string, value string, table string) ([]map[string]any, error)
	Lookup(ctx context.Context, dataType, value, module string, columns []string) (*VtigerResponse[[]map[string]any], error)
	AddRelated(ctx context.Context, source string, related string, label string) (*VtigerResponse[[]map[string]any], error)
	Query(ctx context.Context, query string) (*VtigerResponse[[]map[string]any], error)
	RetrieveRelated(ctx context.Context, id string, module string) (*VtigerResponse[[]map[string]any], error)
	Retrieve(ctx context.Context, id string) (*VtigerResponse[map[string]any], error)
	Describe(ctx context.Context, element string) (*VtigerResponse[Module], error)
	Delete(ctx context.Context, element string) error
	RetrieveFiles(ctx context.Context, id string) (*VtigerResponse[[]File], error)
	Update(ctx context.Context, data map[string]any) (*VtigerResponse[map[string]any], error)
	Revise(ctx context.Context, data map[string]any) (*VtigerResponse[map[string]any], error)
	Create(ctx context.Context, element string, data map[string]any) (*VtigerResponse[map[string]any], error)
	Count(ctx context.Context, module string, filters map[string]string) (int, error)
	ExecuteCount(ctx context.Context, query string) (int, error)
}
