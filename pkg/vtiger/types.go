package vtiger

import "context"

type Creator interface {
	Create(ctx context.Context, element string, data map[string]any) (*VtigerResponse[map[string]any], error)
}
