package vtiger

import (
	"encoding/json"
	"time"
)

type Model struct {
	CreatedTime    time.Time `json:"created_time"`
	ModifiedTime   time.Time `json:"modified_time"`
	AssignedUserId string    `json:"assigned_user_id"`
	Description    string    `json:"description"`
	Source         string    `json:"source"`
	Starred        bool      `json:"starred"`
	Tags           []string  `json:"tags"`
	Id             string    `json:"id"`
	Label          string    `json:"label"`
}

func (m *Model) ConvertToMap() map[string]any {
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}

	return result
}
