package domain

import (
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Servicename      string    `json:"servicename"`
	ServiceNo        string    `json:"service_no"`
	Discontinued     bool      `json:"discontinued"`
	SalesStartDate   time.Time `json:"sales_start_date"`
	SalesEndDate     time.Time `json:"sales_end_date"`
	StartDate        time.Time `json:"start_date"`
	ExpiryDate       time.Time `json:"expiry_date"`
	Website          string    `json:"website"`
	ServiceUsageunit string    `json:"service_usageunit"`
	QtyPerUnit       float64   `json:"qty_per_unit"`
	Servicecategory  string    `json:"servicecategory"`
	UnitPrice        float64   `json:"unit_price"`
	Taxclass         string    `json:"taxclass"`
	Commissionrate   float64   `json:"commissionrate"`
	PurchaseCost     float64   `json:"purchase_cost"`
	Tax2             float64   `json:"tax2"`
	Tax3             float64   `json:"tax3"`
	Currency1        float64   `json:"currency1"`
	CurrencyId       string    `json:"currency_id"`
	CreatedTime      time.Time `json:"created_time"`
	ModifiedTime     time.Time `json:"modified_time"`
	AssignedUserId   string    `json:"assigned_user_id"`
	Description      string    `json:"description"`
	Source           string    `json:"source"`
	Starred          bool      `json:"starred"`
	Tags             []string  `json:"tags"`
	Id               string    `json:"id"`
	Label            string    `json:"label"`
	Currency         Currency  `json:"currency"`
}

var MockedService = Service{
	Servicename:      "Cleaning laptop",
	ServiceNo:        "SER1",
	Discontinued:     true,
	SalesStartDate:   time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	SalesEndDate:     time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	StartDate:        time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	ExpiryDate:       time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	Website:          "https://itvolga.com",
	ServiceUsageunit: "Hours",
	QtyPerUnit:       4.0,
	Servicecategory:  "Support",
	UnitPrice:        15.0,
	Taxclass:         "",
	Commissionrate:   2.5,
	PurchaseCost:     56.4,
	Tax2:             10.0,
	Tax3:             12.5,
	Currency1:        15,
	CurrencyId:       "21x1",
	CreatedTime:      time.Date(2023, time.April, 25, 0, 0, 0, 0, time.UTC),
	ModifiedTime:     time.Date(2023, time.April, 30, 0, 0, 0, 0, time.UTC),
	AssignedUserId:   "19x1",
	Description:      "Some Description",
	Source:           "CRM",
	Starred:          false,
	Tags:             []string{"test1"},
	Id:               "25x52",
	Label:            "Cleaning laptop",
}

func ConvertMapToService(inputMap map[string]any) (Service, error) {
	var service Service
	layout := "2006-01-02 15:04:05"
	shortLayout := "2006-01-02"

	for key, value := range inputMap {
		switch key {
		case "id":
			service.Id = value.(string)
		case "servicename":
			service.Servicename = value.(string)
		case "service_usageunit":
			service.ServiceUsageunit = value.(string)
		case "servicecategory":
			service.Servicecategory = value.(string)
		case "service_no":
			service.ServiceNo = value.(string)
		case "assigned_user_id":
			service.AssignedUserId = value.(string)
		case "start_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			service.StartDate = parsedTime
		case "sales_start_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			service.SalesStartDate = parsedTime
		case "sales_end_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			service.SalesEndDate = parsedTime
		case "expiry_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			service.ExpiryDate = parsedTime
		case "website":
			service.Website = value.(string)
		case "createdtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			service.CreatedTime = parsedTime
		case "modifiedtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			service.ModifiedTime = parsedTime
		case "currency1":
			service.Currency1 = value.(float64)
		case "unit_price", "commissionrate", "qty_per_unit", "purchase_cost", "tax2", "tax3":
			val := value.(string)
			f, err := strconv.ParseFloat(val, 64)
			if err == nil {
				switch key {
				case "unit_price":
					service.UnitPrice = f
				case "commissionrate":
					service.Commissionrate = f
				case "qty_per_unit":
					service.QtyPerUnit = f
				case "purchase_cost":
					service.PurchaseCost = f
				case "tax2":
					service.Tax2 = f
				case "tax3":
					service.Tax3 = f
				}
			}
		case "starred":
			service.Starred = value.(string) == "1"
		case "taxclass":
			service.Taxclass = value.(string)
		case "description":
			service.Description = value.(string)
		case "discontinued":
			service.Discontinued = value.(string) == "1"
		case "label":
			service.Label = value.(string)
		case "currency_id":
			service.CurrencyId = value.(string)
		case "tags":
			service.Tags = strings.Split(value.(string), ",")
		}
	}

	return service, nil
}
