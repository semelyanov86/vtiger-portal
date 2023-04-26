package utils

import (
	"encoding/json"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

func ConvertStructToMap(conv any) (map[string]any, error) {
	var result map[string]interface{}
	data, err := json.Marshal(conv)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, e.Wrap("can not convert struct to map", err)
	}
	return result, nil
}
