package domain

import (
	"github.com/semelyanov86/vtiger-portal/internal/utils"
	"time"
)

type Comment struct {
	Id             string        `json:"id"`
	Commentcontent string        `json:"commentcontent"`
	Source         string        `json:"source"`
	Customer       string        `json:"customer"`
	Userid         string        `json:"userid"`
	Reasontoedit   string        `json:"reasontoedit"`
	Creator        string        `json:"creator"`
	AssignedUserId string        `json:"assigned_user_id"`
	Createdtime    time.Time     `json:"createdtime"`
	Modifiedtime   time.Time     `json:"modifiedtime"`
	RelatedTo      string        `json:"related_to"`
	ParentComments string        `json:"parent_comments"`
	IsPrivate      bool          `json:"is_private"`
	Filename       string        `json:"filename"`
	RelatedEmailId string        `json:"related_email_id"`
	Author         CommentAuthor `json:"author"`
}

type CommentAuthor struct {
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Email        string `json:"email"`
	Imagecontent string `json:"imagecontent"`
	Id           string `json:"id"`
}

var MockedComment = Comment{
	Id:             "37x6625",
	Commentcontent: "This is a test comment.",
	Source:         "CRM",
	Customer:       "",
	Userid:         "",
	Reasontoedit:   "Typo fix",
	Creator:        "19x1",
	AssignedUserId: "19x1",
	Createdtime:    time.Now().Add(-time.Hour),
	Modifiedtime:   time.Now(),
	RelatedTo:      "17x923",
	ParentComments: "",
	IsPrivate:      false,
	Filename:       "",
	RelatedEmailId: "",
}

func ConvertMapToComment(m map[string]any) Comment {
	comment := Comment{}
	layout := "2006-01-02 15:04:05"

	for k, v := range m {
		switch k {
		case "id":
			comment.Id = v.(string)
		case "commentcontent":
			comment.Commentcontent = v.(string)
		case "source":
			comment.Source = v.(string)
		case "customer":
			comment.Customer = v.(string)
		case "userid":
			comment.Userid = v.(string)
		case "reasontoedit":
			comment.Reasontoedit = v.(string)
		case "creator":
			comment.Creator = v.(string)
		case "assigned_user_id":
			comment.AssignedUserId = v.(string)
		case "related_to":
			comment.RelatedTo = v.(string)
		case "createdtime":
			parsedTime, _ := time.Parse(layout, v.(string))
			comment.Createdtime = parsedTime
		case "modifiedtime":
			parsedTime, _ := time.Parse(layout, v.(string))
			comment.Modifiedtime = parsedTime
		case "parent_comments":
			comment.ParentComments = v.(string)
		case "is_private":
			comment.IsPrivate = v.(string) == "1"
		case "filename":
			comment.Filename = v.(string)
		case "related_email_id":
			comment.RelatedEmailId = v.(string)
		}
	}

	return comment
}

func (c Comment) ConvertToMap() (map[string]any, error) {
	return utils.ConvertStructToMap(c)
}
