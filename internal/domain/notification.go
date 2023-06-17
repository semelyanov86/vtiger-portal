package domain

import "time"

type Notification struct {
	Id             int64     `json:"id"`
	Crmid          string    `json:"crmid"`
	Module         string    `json:"module"`
	Label          string    `json:"label"`
	Description    string    `json:"description"`
	Manager        Manager   `json:"manager"`
	AssignedUserId string    `json:"-"`
	AccountId      string    `json:"-"`
	UserId         string    `json:"-"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updatead_at"`
}

func ConvertMapToNotification(data map[string]any) (Notification, error) {
	notification := Notification{}

	if crmid, ok := data["crmid"].(string); ok {
		notification.Crmid = crmid
	}
	if module, ok := data["module"].(string); ok {
		notification.Module = module
	}
	if label, ok := data["label"].(string); ok {
		notification.Label = label
	}
	if description, ok := data["description"].(string); ok {
		notification.Description = description
	}
	if manager, ok := data["manager"].(Manager); ok {
		notification.Manager = manager
	}
	if isRead, ok := data["is_read"]; ok {
		if isRead == "1" || isRead == 1 {
			notification.IsRead = true
		} else {
			notification.IsRead = false
		}
	}

	return notification, nil
}
