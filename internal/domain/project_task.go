package domain

import (
	"strings"
	"time"
)

type ProjectTask struct {
	Projecttaskname     string    `json:"projecttaskname"`
	Projecttasktype     string    `json:"projecttasktype"`
	Projecttaskpriority string    `json:"projecttaskpriority"`
	Projectid           string    `json:"projectid"`
	AssignedUserId      string    `json:"assigned_user_id"`
	Projecttasknumber   string    `json:"projecttasknumber"`
	ProjecttaskNo       string    `json:"projecttask_no"`
	Projecttaskprogress string    `json:"projecttaskprogress"`
	Projecttaskhours    string    `json:"projecttaskhours"`
	Startdate           time.Time `json:"startdate"`
	Enddate             time.Time `json:"enddate"`
	CreatedTime         time.Time `json:"createdtime"`
	ModifiedTime        time.Time `json:"modifiedtime"`
	Description         string    `json:"description"`
	Source              string    `json:"source"`
	Starred             bool      `json:"starred"`
	Tags                []string  `json:"tags"`
	Projecttaskstatus   string    `json:"projecttaskstatus"`
	Id                  string    `json:"id"`
	Label               string    `json:"label"`
}

var MockedProjectTask = ProjectTask{
	Projecttaskname:     "Install hosting",
	Projecttasktype:     "administrative",
	Projecttaskpriority: "low",
	Projectid:           "29x54",
	AssignedUserId:      "19x1",
	Projecttasknumber:   "2",
	ProjecttaskNo:       "PT1",
	Projecttaskprogress: "20%",
	Projecttaskhours:    "15",
	Startdate:           time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	Enddate:             time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	CreatedTime:         time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	ModifiedTime:        time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	Description:         "Some descr",
	Source:              "CRM",
	Starred:             false,
	Tags:                []string{"test"},
	Projecttaskstatus:   "Open",
	Id:                  "28x56",
	Label:               "Install hosting",
}

func ConvertMapToProjectTask(inputMap map[string]any) (ProjectTask, error) {
	var projectTask ProjectTask
	layout := "2006-01-02 15:04:05"
	shortLayout := "2006-01-02"

	for key, value := range inputMap {
		switch key {
		case "id":
			projectTask.Id = value.(string)
		case "projecttaskname":
			projectTask.Projecttaskname = value.(string)
		case "projecttasktype":
			projectTask.Projecttasktype = value.(string)
		case "projecttaskpriority":
			projectTask.Projecttaskpriority = value.(string)
		case "projectid":
			projectTask.Projectid = value.(string)
		case "assigned_user_id":
			projectTask.AssignedUserId = value.(string)
		case "startdate":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			projectTask.Startdate = parsedTime
		case "enddate":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			projectTask.Enddate = parsedTime
		case "projecttask_no":
			projectTask.ProjecttaskNo = value.(string)
		case "projecttaskprogress":
			projectTask.Projecttaskprogress = value.(string)
		case "createdtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			projectTask.CreatedTime = parsedTime
		case "modifiedtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			projectTask.ModifiedTime = parsedTime
		case "projecttaskstatus":
			projectTask.Projecttaskstatus = value.(string)
		case "projecttasknumber":
			projectTask.Projecttasknumber = value.(string)
		case "starred":
			projectTask.Starred = value.(string) == "1"
		case "projecttaskhours":
			projectTask.Projecttaskhours = value.(string)
		case "description":
			projectTask.Description = value.(string)
		case "label":
			projectTask.Label = value.(string)
		case "source":
			projectTask.Source = value.(string)
		case "tags":
			projectTask.Tags = strings.Split(value.(string), ",")
		}
	}

	return projectTask, nil
}
