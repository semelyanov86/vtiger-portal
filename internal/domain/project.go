package domain

import (
	"strings"
	"time"
)

type Project struct {
	Projectname              string    `json:"projectname"`
	Startdate                time.Time `json:"startdate"`
	Targetenddate            time.Time `json:"targetenddate"`
	Actualenddate            time.Time `json:"actualenddate"`
	Projectstatus            string    `json:"projectstatus"`
	Projecttype              string    `json:"projecttype"`
	Linktoaccountscontacts   string    `json:"linktoaccountscontacts"`
	ProjectNo                string    `json:"project_no"`
	Targetbudget             string    `json:"targetbudget"`
	Projecturl               string    `json:"projecturl"`
	Projectpriority          string    `json:"projectpriority"`
	Progress                 string    `json:"progress"`
	Isconvertedfrompotential bool      `json:"isconvertedfrompotential"`
	Potentialid              string    `json:"potentialid"`
	CreatedTime              time.Time `json:"created_time"`
	ModifiedTime             time.Time `json:"modified_time"`
	AssignedUserId           string    `json:"assigned_user_id"`
	Description              string    `json:"description"`
	Source                   string    `json:"source"`
	Starred                  bool      `json:"starred"`
	Tags                     []string  `json:"tags"`
	Id                       string    `json:"id"`
	Label                    string    `json:"label"`
}

var MockedProject = Project{
	Projectname:              "Website development",
	Startdate:                time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	Targetenddate:            time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	Actualenddate:            time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	Projectstatus:            "prospecting",
	Projecttype:              "administrative",
	Linktoaccountscontacts:   "11x1",
	ProjectNo:                "PROJ1",
	Targetbudget:             "500 RUB",
	Projecturl:               "https://sergeyem.ru",
	Projectpriority:          "normal",
	Progress:                 "10%",
	Isconvertedfrompotential: false,
	Potentialid:              "13x14",
	CreatedTime:              time.Date(2023, time.April, 25, 0, 0, 0, 0, time.UTC),
	ModifiedTime:             time.Date(2023, time.April, 25, 0, 0, 0, 0, time.UTC),
	AssignedUserId:           "19x1",
	Description:              "This is test project",
	Source:                   "CRM",
	Starred:                  false,
	Tags:                     nil,
	Id:                       "29x54",
	Label:                    "Website development",
}

func ConvertMapToProject(inputMap map[string]any) (Project, error) {
	var project Project
	layout := "2006-01-02 15:04:05"
	shortLayout := "2006-01-02"

	for key, value := range inputMap {
		switch key {
		case "id":
			project.Id = value.(string)
		case "projectname":
			project.Projectname = value.(string)
		case "projectstatus":
			project.Projectstatus = value.(string)
		case "projecttype":
			project.Projecttype = value.(string)
		case "linktoaccountscontacts":
			project.Linktoaccountscontacts = value.(string)
		case "assigned_user_id":
			project.AssignedUserId = value.(string)
		case "startdate":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			project.Startdate = parsedTime
		case "targetenddate":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			project.Targetenddate = parsedTime
		case "actualenddate":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			project.Actualenddate = parsedTime
		case "project_no":
			project.ProjectNo = value.(string)
		case "progress":
			project.Progress = value.(string)
		case "createdtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			project.CreatedTime = parsedTime
		case "modifiedtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			project.ModifiedTime = parsedTime
		case "targetbudget":
			project.Targetbudget = value.(string)
		case "potentialid":
			project.Potentialid = value.(string)
		case "starred":
			project.Starred = value.(string) == "1"
		case "isconvertedfrompotential":
			project.Isconvertedfrompotential = value.(string) == "1"
		case "projecturl":
			project.Projecturl = value.(string)
		case "description":
			project.Description = value.(string)
		case "label":
			project.Label = value.(string)
		case "projectpriority":
			project.Projectpriority = value.(string)
		case "tags":
			project.Tags = strings.Split(value.(string), ",")
		}
	}

	return project, nil
}
