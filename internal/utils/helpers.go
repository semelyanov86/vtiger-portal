package utils

import (
	"encoding/json"
	"errors"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
	"time"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

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

func BytesToDate(jsonValue []byte, layout string) (time.Time, error) {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return time.Now(), ErrInvalidRuntimeFormat
	}
	parsedTime, err := time.Parse(layout, unquotedJSONValue)
	if err != nil {
		return parsedTime, e.Wrap("can not parse date "+unquotedJSONValue, err)
	}
	return parsedTime, nil
}

func BytesToBool(jsonValue []byte) (bool, error) {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return false, ErrInvalidRuntimeFormat
	}

	return unquotedJSONValue == "1", nil
}

func FindFieldByRefers(module vtiger.Module, refers string) *vtiger.ModuleField {
	for _, field := range module.Fields {
		if len(field.Type.RefersTo) > 0 {
			for _, s := range field.Type.RefersTo {
				if s == refers {
					return &field
				}
			}
		}
	}
	return nil
}
