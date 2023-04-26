package domain

import (
	"github.com/semelyanov86/vtiger-portal/internal/utils"
	"strings"
	"time"
)

type Faq struct {
	ProductId     string    `json:"product_id"`
	ID            string    `json:"id"`
	Faqcategories string    `json:"faqcategories"`
	Faqstatus     string    `json:"faqstatus"`
	Question      string    `json:"question"`
	FaqAnswer     string    `json:"faq_answer"`
	CreatedTime   time.Time `json:"createdtime"`
	ModifiedTime  time.Time `json:"modifiedtime"`
	Starred       bool      `json:"starred"`
	Tags          []string  `json:"tags"`
}

var MockedFaq = Faq{
	ProductId:     "16x5",
	ID:            "6x50",
	Faqcategories: "General",
	Faqstatus:     "Published",
	Question:      "How to write text?",
	FaqAnswer:     "Just write it and that is it",
	CreatedTime:   time.Now().Add(-24 * time.Hour), // Example: created 24 hours ago
	ModifiedTime:  time.Now().Add(-12 * time.Hour), // Example: modified 12 hours ago
	Starred:       false,
	Tags:          []string{"mock_tag1", "mock_tag2"},
}

func ConvertMapToFaq(m map[string]any) (Faq, error) {
	faq := Faq{}
	layout := "2006-01-02 15:04:05"

	for k, v := range m {
		switch k {
		case "id":
			faq.ID = v.(string)
		case "product_id":
			faq.ProductId = v.(string)
		case "faqcategories":
			faq.Faqcategories = v.(string)
		case "faqstatus":
			faq.Faqstatus = v.(string)
		case "question":
			faq.Question = v.(string)
		case "faq_answer":
			faq.FaqAnswer = v.(string)
		case "createdtime":
			parsedTime, _ := time.Parse(layout, v.(string))
			faq.CreatedTime = parsedTime
		case "modifiedtime":
			parsedTime, _ := time.Parse(layout, v.(string))
			faq.ModifiedTime = parsedTime
		case "starred":
			faq.Starred = v.(string) == "1"
		case "tags":
			faq.Tags = strings.Split(v.(string), ",")
		}
	}

	return faq, nil
}

func (f Faq) ConvertToMap() (map[string]any, error) {
	result, err := utils.ConvertStructToMap(f)
	if err != nil {
		return result, err
	}
	tags := ""
	if len(f.Tags) > 0 {
		tags = strings.Join(f.Tags, ",")
	}
	result["tags"] = tags
	return result, nil
}
