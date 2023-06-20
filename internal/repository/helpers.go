package repository

func GetIsActiveFromFilter(filters map[string]any) string {
	isActive := "1"
	if filters["discontinued"] == false || filters["discontinued"] == "false" {
		isActive = "0"
	}
	return isActive
}
