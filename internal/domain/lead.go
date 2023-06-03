package domain

import (
	"github.com/semelyanov86/vtiger-portal/internal/utils"
)

type Lead struct {
	ID             string `json:"id"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname" binding:"required"`
	Phone          string `json:"phone" binding:"required"`
	Company        string `json:"company"`
	Email          string `json:"email" binding:"required"`
	Mobile         string `json:"mobile"`
	Website        string `json:"website"`
	Leadstatus     string `json:"leadstatus"`
	Leadsource     string `json:"leadsource"`
	AssignedUserId string `json:"assigned_user_id"`
	Description    string `json:"description"`
}

func (c Lead) ConvertToMap() (map[string]any, error) {
	return utils.ConvertStructToMap(c)
}

func ConvertMapToLead(m map[string]any) Lead {
	lead := Lead{}

	for k, v := range m {
		switch k {
		case "id":
			lead.ID = v.(string)
		case "firstname":
			lead.Firstname = v.(string)
		case "lastname":
			lead.Lastname = v.(string)
		case "phone":
			lead.Phone = v.(string)
		case "company":
			lead.Company = v.(string)
		case "email":
			lead.Email = v.(string)
		case "mobile":
			lead.Mobile = v.(string)
		case "assigned_user_id":
			lead.AssignedUserId = v.(string)
		case "website":
			lead.Website = v.(string)
		case "leadstatus":
			lead.Leadstatus = v.(string)
		case "leadsource":
			lead.Leadsource = v.(string)
		case "description":
			lead.Description = v.(string)
		}
	}

	return lead
}
