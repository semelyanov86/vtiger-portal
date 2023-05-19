package domain

import "time"

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
