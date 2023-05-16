package domain

import (
	"strings"
	"time"
)

type ServiceContract struct {
	ID               string    `json:"id"`
	AssignedUserId   string    `json:"assigned_user_id"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	ScRelatedTo      string    `json:"sc_related_to"`
	TrackingUnit     string    `json:"tracking_unit"`
	TotalUnits       float64   `json:"total_units"`
	UsedUnits        float64   `json:"used_units"`
	Subject          string    `json:"subject"`
	DueDate          time.Time `json:"due_date"`
	PlannedDuration  int       `json:"planned_duration"`
	ActualDuration   int       `json:"actual_duration"`
	ContractStatus   string    `json:"contract_status"`
	ContractPriority string    `json:"contract_priority"`
	ContractType     string    `json:"contract_type"`
	Progress         float64   `json:"progress"`
	ContractNo       string    `json:"contract_no"`
	CreatedTime      time.Time `json:"createdtime"`
	ModifiedTime     time.Time `json:"modifiedtime"`
	Starred          bool      `json:"starred"`
	Tags             []string  `json:"tags"`
	Label            string    `json:"label"`
}

var MockedServiceContract = ServiceContract{
	ID:               "24x59",
	AssignedUserId:   "19x1",
	StartDate:        time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	EndDate:          time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
	ScRelatedTo:      "12x11",
	TrackingUnit:     "hours",
	TotalUnits:       1000,
	UsedUnits:        200,
	Subject:          "Mocked Service Contract",
	DueDate:          time.Date(2023, time.August, 15, 0, 0, 0, 0, time.UTC),
	PlannedDuration:  120,
	ActualDuration:   90,
	ContractStatus:   "In Progress",
	ContractPriority: "high",
	ContractType:     "Support",
	Progress:         20,
	ContractNo:       "SERCON2",
	CreatedTime:      time.Date(2023, time.April, 25, 0, 0, 0, 0, time.UTC),
	ModifiedTime:     time.Date(2023, time.April, 30, 0, 0, 0, 0, time.UTC),
	Starred:          false,
	Tags:             []string{"tag1", "tag2"},
	Label:            "Label-A",
}

func ConvertMapToServiceContract(inputMap map[string]any) (ServiceContract, error) {
	var serviceContract ServiceContract
	layout := "2006-01-02 15:04:05"
	shortLayout := "2006-01-02"

	for key, value := range inputMap {
		switch key {
		case "id":
			serviceContract.ID = value.(string)
		case "assigned_user_id":
			serviceContract.AssignedUserId = value.(string)
		case "start_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			serviceContract.StartDate = parsedTime
		case "end_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			serviceContract.EndDate = parsedTime
		case "due_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			serviceContract.DueDate = parsedTime
		case "createdtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			serviceContract.CreatedTime = parsedTime
		case "modifiedtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			serviceContract.ModifiedTime = parsedTime
		case "sc_related_to":
			serviceContract.ScRelatedTo = value.(string)
		case "tracking_unit":
			serviceContract.TrackingUnit = value.(string)
		case "total_units", "used_units", "progress":
			if f, ok := value.(float64); ok {
				switch key {
				case "total_units":
					serviceContract.TotalUnits = f
				case "used_units":
					serviceContract.UsedUnits = f
				case "progress":
					serviceContract.Progress = f
				}
			}
		case "subject":
			serviceContract.Subject = value.(string)
		case "planned_duration", "actual_duration":
			if i, ok := value.(int); ok {
				switch key {
				case "planned_duration":
					serviceContract.PlannedDuration = i
				case "actual_duration":
					serviceContract.ActualDuration = i
				}
			}
		case "contract_status":
			serviceContract.ContractStatus = value.(string)
		case "contract_priority":
			serviceContract.ContractPriority = value.(string)
		case "contract_type":
			serviceContract.ContractType = value.(string)
		case "contract_no":
			serviceContract.ContractNo = value.(string)
		case "starred":
			serviceContract.Starred = value.(string) == "1"
		case "tags":
			serviceContract.Tags = strings.Split(value.(string), ",")
		case "label":
			serviceContract.Label = value.(string)
		}
	}

	return serviceContract, nil
}
