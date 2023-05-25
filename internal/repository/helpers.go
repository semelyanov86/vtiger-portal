package repository

import (
	"fmt"
	"strings"
)

func GetIsActiveFromFilter(filters map[string]any) string {
	isActive := "1"
	if filters["discontinued"] == false || filters["discontinued"] == "false" {
		isActive = "0"
	}
	return isActive
}

func GenerateOrderByClause(sort string) string {
	if sort == "" {
		return ""
	}

	fields := strings.Split(sort, ",")
	orderClauses := make([]string, 0, len(fields))

	for _, field := range fields {
		var orderBy string
		if strings.HasPrefix(field, "-") {
			orderBy = fmt.Sprintf("%s DESC", strings.TrimPrefix(field, "-"))
		} else {
			orderBy = fmt.Sprintf("%s ASC", field)
		}
		orderClauses = append(orderClauses, orderBy)
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(orderClauses, ", "))
}
