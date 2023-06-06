package domain

import (
	"github.com/semelyanov86/vtiger-portal/internal/utils"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strings"
	"time"
)

type Document struct {
	NotesTitle         string    `json:"notes_title"`
	Createdtime        time.Time `json:"createdtime"`
	Modifiedtime       time.Time `json:"modifiedtime"`
	Filename           string    `json:"filename"`
	AssignedUserId     string    `json:"assigned_user_id"`
	Notecontent        string    `json:"notecontent"`
	Filetype           string    `json:"filetype"`
	Filesize           string    `json:"filesize"`
	Filelocationtype   string    `json:"filelocationtype"`
	Fileversion        string    `json:"fileversion"`
	Filestatus         string    `json:"filestatus"`
	Filedownloadcount  string    `json:"filedownloadcount"`
	Folderid           string    `json:"folderid"`
	NoteNo             string    `json:"note_no"`
	Modifiedby         string    `json:"modifiedby"`
	Source             string    `json:"source"`
	Starred            bool      `json:"starred"`
	Tags               []string  `json:"tags"`
	Id                 string    `json:"id"`
	Imageattachmentids string    `json:"imageattachmentids"`
}

func ConvertMapToDocument(m map[string]any) Document {
	document := Document{}
	layout := "2006-01-02 15:04:05"

	for k, v := range m {
		switch k {
		case "id":
			document.Id = v.(string)
		case "notes_title":
			document.NotesTitle = v.(string)
		case "source":
			document.Source = v.(string)
		case "notecontent":
			document.Notecontent = v.(string)
		case "filetype":
			document.Filetype = v.(string)
		case "filesize":
			document.Filesize = v.(string)
		case "filelocationtype":
			document.Filelocationtype = v.(string)
		case "assigned_user_id":
			document.AssignedUserId = v.(string)
		case "fileversion":
			document.Fileversion = v.(string)
		case "createdtime":
			parsedTime, _ := time.Parse(layout, v.(string))
			document.Createdtime = parsedTime
		case "modifiedtime":
			parsedTime, _ := time.Parse(layout, v.(string))
			document.Modifiedtime = parsedTime
		case "filestatus":
			document.Filestatus = v.(string)
		case "filedownloadcount":
			document.Filedownloadcount = v.(string)
		case "filename":
			document.Filename = v.(string)
		case "folderid":
			document.Folderid = v.(string)
		case "note_no":
			document.NoteNo = v.(string)
		case "modifiedby":
			document.Modifiedby = v.(string)
		case "imageattachmentids":
			document.Imageattachmentids = v.(string)
		case "starred":
			document.Starred = v.(string) == "1"
		case "tags":
			document.Tags = strings.Split(v.(string), ",")
		}
	}

	return document
}

var MockedDocument = Document{
	NotesTitle:         "customer-portal",
	Createdtime:        time.Now(),
	Modifiedtime:       time.Now(),
	Filename:           "portal.yaml",
	AssignedUserId:     "19x1",
	Notecontent:        "",
	Filetype:           "application/yaml",
	Filesize:           "1277",
	Filelocationtype:   "I",
	Fileversion:        "",
	Filestatus:         "1",
	Filedownloadcount:  "",
	Folderid:           "22x1",
	NoteNo:             "DOC1",
	Modifiedby:         "19x1",
	Source:             "CRM",
	Starred:            false,
	Tags:               []string{"test1"},
	Id:                 "15x40",
	Imageattachmentids: "15x42",
}

var MockedFile = vtiger.File{
	Fileid:       "42",
	Filename:     "async-vtiger.png",
	Filetype:     "image/png",
	Filesize:     19951,
	Filecontents: "iVBORw0KGgoAAAANSUhEUgAAAnAAAAD1CAIAAA",
}

func (d Document) ConvertToMap() (map[string]any, error) {
	return utils.ConvertStructToMap(d)
}
