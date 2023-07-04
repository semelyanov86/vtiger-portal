package vtiger

import (
	"context"
	"time"
)

type MockedConnector struct {
}

var MockedEntity = Model{
	CreatedTime:    time.Now(),
	ModifiedTime:   time.Now(),
	AssignedUserId: "19x1",
	Description:    "Some description for mocked entity",
	Source:         "CRM",
	Starred:        false,
	Tags:           []string{"tag1", "tag2"},
	Id:             "23x42343",
	Label:          "This is test label",
}

func (m MockedConnector) GetAll(ctx context.Context, filter PaginationQueryFilter, fields QueryFieldsProps) ([]map[string]any, error) {
	result := make([]map[string]any, 2)
	result[0] = MockedEntity.ConvertToMap()
	result[1] = MockedEntity.ConvertToMap()
	return result, nil
}

func (m MockedConnector) GetByWhereClause(ctx context.Context, filter PaginationQueryFilter, field string, value string, table string) ([]map[string]any, error) {
	result := make([]map[string]any, 2)
	result[0] = MockedEntity.ConvertToMap()
	result[1] = MockedEntity.ConvertToMap()
	return result, nil
}

func (m MockedConnector) Lookup(ctx context.Context, dataType, value, module string, columns []string) (*VtigerResponse[[]map[string]any], error) {
	result := make([]map[string]any, 2)
	result[0] = MockedEntity.ConvertToMap()
	result[1] = MockedEntity.ConvertToMap()
	return &VtigerResponse[[]map[string]any]{Result: result}, nil
}

func (m MockedConnector) AddRelated(ctx context.Context, source string, related string, label string) (*VtigerResponse[[]map[string]any], error) {
	return nil, nil
}

func (m MockedConnector) Query(ctx context.Context, query string) (*VtigerResponse[[]map[string]any], error) {
	result := make([]map[string]any, 2)
	result[0] = MockedEntity.ConvertToMap()
	result[1] = MockedEntity.ConvertToMap()
	return &VtigerResponse[[]map[string]any]{Result: result}, nil
}

func (m MockedConnector) RetrieveRelated(ctx context.Context, id string, module string) (*VtigerResponse[[]map[string]any], error) {
	result := make([]map[string]any, 2)
	result[0] = MockedEntity.ConvertToMap()
	result[1] = MockedEntity.ConvertToMap()
	return &VtigerResponse[[]map[string]any]{Result: result}, nil
}

func (m MockedConnector) Retrieve(ctx context.Context, id string) (*VtigerResponse[map[string]any], error) {
	MockedEntity.Id = id
	return &VtigerResponse[map[string]any]{Result: MockedEntity.ConvertToMap()}, nil
}

func (m MockedConnector) Describe(ctx context.Context, element string) (*VtigerResponse[Module], error) {
	return &VtigerResponse[Module]{Result: MockedModule, Success: true}, nil
}

func (m MockedConnector) Delete(ctx context.Context, element string) error {
	return nil
}

func (m MockedConnector) RetrieveFiles(ctx context.Context, id string) (*VtigerResponse[[]File], error) {
	return nil, nil
}

func (m MockedConnector) Update(ctx context.Context, data map[string]any) (*VtigerResponse[map[string]any], error) {
	return nil, nil
}

func (m MockedConnector) Revise(ctx context.Context, data map[string]any) (*VtigerResponse[map[string]any], error) {
	return nil, nil
}

func (m MockedConnector) Create(ctx context.Context, element string, data map[string]any) (*VtigerResponse[map[string]any], error) {
	return &VtigerResponse[map[string]any]{Result: MockedEntity.ConvertToMap()}, nil
}

func (m MockedConnector) Count(ctx context.Context, module string, filters map[string]string) (int, error) {
	return 2, nil
}

func (m MockedConnector) ExecuteCount(ctx context.Context, query string) (int, error) {
	return 2, nil
}

func NewMockedVtigerConnector() *MockedConnector {
	return &MockedConnector{}
}
