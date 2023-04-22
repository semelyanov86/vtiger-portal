package domain

import "time"

type Comment struct {
	Id             string    `json:"id"`
	Commentcontent string    `json:"commentcontent"`
	Source         string    `json:"source"`
	Customer       string    `json:"customer"`
	Userid         string    `json:"userid"`
	Reasontoedit   string    `json:"reasontoedit"`
	Creator        string    `json:"creator"`
	AssignedUserId string    `json:"assigned_user_id"`
	Createdtime    time.Time `json:"createdtime"`
	Modifiedtime   time.Time `json:"modifiedtime"`
	RelatedTo      string    `json:"related_to"`
	ParentComments string    `json:"parent_comments"`
	IsPrivate      bool      `json:"is_private"`
	Filename       string    `json:"filename"`
	RelatedEmailId string    `json:"related_email_id"`
}
