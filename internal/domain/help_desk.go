package domain

import (
	"errors"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"strconv"
	"strings"
	"time"
)

var ErrCanNotConvertValue = errors.New("can not convert value")

type HelpDesk struct {
	TicketNo         string    `json:"ticket_no"`
	AssignedUserID   string    `json:"assigned_user_id"`
	ParentID         string    `json:"parent_id"`
	TicketPriorities string    `json:"ticketpriorities"`
	ProductID        string    `json:"product_id"`
	TicketSeverities string    `json:"ticketseverities"`
	TicketStatus     string    `json:"ticketstatus"`
	TicketCategories string    `json:"ticketcategories"`
	Hours            float64   `json:"hours"`
	Days             float64   `json:"days"`
	CreatedTime      time.Time `json:"createdtime"`
	ModifiedTime     time.Time `json:"modifiedtime"`
	FromPortal       bool      `json:"from_portal"`
	ModifiedBy       string    `json:"modifiedby"`
	TicketTitle      string    `json:"ticket_title"`
	Description      string    `json:"description"`
	Solution         string    `json:"solution"`
	ContactID        string    `json:"contact_id"`
	CreatedUserID    string    `json:"created_user_id"`
	Source           string    `json:"source"`
	Starred          bool      `json:"starred"`
	Tags             []string  `json:"tags"`
	ID               string    `json:"id"`
	Label            string    `json:"label"`
}

var MockedHelpDesk = HelpDesk{
	TicketNo:         "TICKET_28",
	AssignedUserID:   "19x1",
	ParentID:         "11x1",
	TicketPriorities: "Normal",
	ProductID:        "",
	TicketSeverities: "",
	TicketStatus:     "Open",
	TicketCategories: "",
	Hours:            10,
	Days:             1,
	CreatedTime:      time.Date(2017, 2, 8, 10, 56, 7, 0, time.UTC),
	ModifiedTime:     time.Date(2018, 7, 13, 18, 55, 51, 0, time.UTC),
	FromPortal:       false,
	ModifiedBy:       "19x1",
	TicketTitle:      "Problem with emails",
	Description:      "They are not attached to client",
	Solution:         "Solution not provided yet",
	ContactID:        "",
	CreatedUserID:    "19x1",
	Source:           "IMPORT",
	Starred:          false,
	Tags:             []string{},
	ID:               "17x923",
	Label:            "Problem with emails",
}

func ConvertMapToHelpDesk(m map[string]any) (HelpDesk, error) {
	helpDesk := HelpDesk{}
	layout := "2006-01-02 15:04:05"

	for k, v := range m {
		switch k {
		case "id":
			helpDesk.ID = v.(string)
		case "ticket_no":
			helpDesk.TicketNo = v.(string)
		case "assigned_user_id":
			helpDesk.AssignedUserID = v.(string)
		case "parent_id":
			helpDesk.ParentID = v.(string)
		case "ticketpriorities":
			helpDesk.TicketPriorities = v.(string)
		case "product_id":
			helpDesk.ProductID = v.(string)
		case "ticketseverities":
			helpDesk.TicketSeverities = v.(string)
		case "ticketstatus":
			helpDesk.TicketStatus = v.(string)
		case "ticketcategories":
			helpDesk.TicketCategories = v.(string)
		case "hours":
			value, err := strconv.ParseFloat(v.(string), 64)
			if err != nil {
				return HelpDesk{}, e.Wrap("can not convert hours "+v.(string), ErrCanNotConvertValue)
			}
			helpDesk.Hours = value
		case "days":
			value, err := strconv.ParseFloat(v.(string), 64)
			if err != nil {
				return HelpDesk{}, e.Wrap("can not convert days "+v.(string), ErrCanNotConvertValue)
			}
			helpDesk.Days = value
		case "createdtime":
			parsedTime, _ := time.Parse(layout, v.(string))
			helpDesk.CreatedTime = parsedTime
		case "modifiedtime":
			parsedTime, _ := time.Parse(layout, v.(string))
			helpDesk.ModifiedTime = parsedTime
		case "from_portal":
			helpDesk.FromPortal = v.(string) == "1"
		case "modifiedby":
			helpDesk.ModifiedBy = v.(string)
		case "ticket_title":
			helpDesk.TicketTitle = v.(string)
		case "description":
			helpDesk.Description = v.(string)
		case "solution":
			helpDesk.Solution = v.(string)
		case "contact_id":
			helpDesk.ContactID = v.(string)
		case "created_user_id":
			helpDesk.CreatedUserID = v.(string)
		case "source":
			helpDesk.Source = v.(string)
		case "starred":
			helpDesk.Starred = v.(string) == "1"
		case "tags":
			helpDesk.Tags = strings.Split(v.(string), ",")
		case "label":
			helpDesk.Label = v.(string)
		}
	}

	return helpDesk, nil
}
