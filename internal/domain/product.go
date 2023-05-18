package domain

import (
	"strconv"
	"time"
)

type Product struct {
	Productname        string    `json:"productname"`
	ProductNo          string    `json:"product_no"`
	Productcode        string    `json:"productcode"`
	Discontinued       bool      `json:"discontinued"`
	Manufacturer       string    `json:"manufacturer"`
	Productcategory    string    `json:"productcategory"`
	SalesStartDate     time.Time `json:"sales_start_date"`
	SalesEndDate       time.Time `json:"sales_end_date"`
	StartDate          time.Time `json:"start_date"`
	ExpiryDate         time.Time `json:"expiry_date"`
	Website            string    `json:"website"`
	VendorId           string    `json:"vendor_id"`
	MfrPartNo          string    `json:"mfr_part_no"`
	VendorPartNo       string    `json:"vendor_part_no"`
	SerialNo           string    `json:"serial_no"`
	Productsheet       string    `json:"productsheet"`
	Glacct             string    `json:"glacct"`
	Createdtime        time.Time `json:"createdtime"`
	Modifiedtime       time.Time `json:"modifiedtime"`
	UnitPrice          float64   `json:"unit_price"`
	Commissionrate     float64   `json:"commissionrate"`
	Taxclass           string    `json:"taxclass"`
	Usageunit          string    `json:"usageunit"`
	QtyPerUnit         float64   `json:"qty_per_unit"`
	Qtyinstock         float64   `json:"qtyinstock"`
	Reorderlevel       int       `json:"reorderlevel"`
	AssignedUserId     string    `json:"assigned_user_id"`
	Qtyindemand        int       `json:"qtyindemand"`
	Description        string    `json:"description"`
	PurchaseCost       float64   `json:"purchase_cost"`
	Starred            bool      `json:"starred"`
	Id                 string    `json:"id"`
	Imageattachmentids string    `json:"imageattachmentids"`
	Label              string    `json:"label"`
	Currency1          float64   `json:"currency1"`
	CurrencyId         string    `json:"currency_id"`
	Currency           Currency  `json:"currency"`
	Imagecontent       string    `json:"imagecontent"`
}

var MockedProduct = Product{
	Productname:        "Keyboard Logitech",
	ProductNo:          "PRO1",
	Productcode:        "KEY-112",
	Discontinued:       false,
	Manufacturer:       "MetBeat Corp",
	Productcategory:    "Hardware",
	SalesStartDate:     time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
	SalesEndDate:       time.Date(2023, time.May, 4, 0, 0, 0, 0, time.UTC),
	StartDate:          time.Date(2023, time.May, 5, 0, 0, 0, 0, time.UTC),
	ExpiryDate:         time.Date(2023, time.May, 25, 0, 0, 0, 0, time.UTC),
	Website:            "https://itvolga.com",
	VendorId:           "",
	MfrPartNo:          "224-ffre",
	VendorPartNo:       "22-RRE",
	SerialNo:           "342345-342345325-324",
	Productsheet:       "Some Sheet",
	Glacct:             "302-Rental-Income",
	Createdtime:        time.Date(2023, time.April, 25, 0, 0, 0, 0, time.UTC),
	Modifiedtime:       time.Date(2023, time.April, 30, 0, 0, 0, 0, time.UTC),
	UnitPrice:          50.0,
	Commissionrate:     2.0,
	Taxclass:           "",
	Usageunit:          "Each",
	QtyPerUnit:         6.0,
	Qtyinstock:         100.0,
	Reorderlevel:       50,
	AssignedUserId:     "19x1",
	Qtyindemand:        30,
	Description:        "Some description",
	PurchaseCost:       500.0,
	Starred:            false,
	Id:                 "14x9",
	Imageattachmentids: "14x62",
	Label:              "Keyboard Logitech",
	Currency1:          50,
	CurrencyId:         "21x11",
	Currency:           MockedCurrency,
	Imagecontent:       "",
}

func ConvertMapToProduct(inputMap map[string]any) (Product, error) {
	var product Product
	layout := "2006-01-02 15:04:05"
	shortLayout := "2006-01-02"

	for key, value := range inputMap {
		switch key {
		case "id":
			product.Id = value.(string)
		case "productname":
			product.Productname = value.(string)
		case "product_no":
			product.ProductNo = value.(string)
		case "productcode":
			product.Productcode = value.(string)
		case "manufacturer":
			product.Manufacturer = value.(string)
		case "productcategory":
			product.Productcategory = value.(string)
		case "assigned_user_id":
			product.AssignedUserId = value.(string)
		case "start_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			product.StartDate = parsedTime
		case "sales_start_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			product.SalesStartDate = parsedTime
		case "sales_end_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			product.SalesEndDate = parsedTime
		case "expiry_date":
			parsedTime, _ := time.Parse(shortLayout, value.(string))
			product.ExpiryDate = parsedTime
		case "website":
			product.Website = value.(string)
		case "vendor_id":
			product.VendorId = value.(string)
		case "mfr_part_no":
			product.MfrPartNo = value.(string)
		case "vendor_part_no":
			product.VendorPartNo = value.(string)
		case "serial_no":
			product.SerialNo = value.(string)
		case "productsheet":
			product.Productsheet = value.(string)
		case "glacct":
			product.Glacct = value.(string)
		case "createdtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			product.Createdtime = parsedTime
		case "modifiedtime":
			parsedTime, _ := time.Parse(layout, value.(string))
			product.Modifiedtime = parsedTime
		case "currency1":
			product.Currency1 = value.(float64)
		case "unit_price", "commissionrate", "qty_per_unit", "qtyinstock", "purchase_cost":
			val := value.(string)
			f, err := strconv.ParseFloat(val, 64)
			if err == nil {
				switch key {
				case "unit_price":
					product.UnitPrice = f
				case "commissionrate":
					product.Commissionrate = f
				case "qty_per_unit":
					product.QtyPerUnit = f
				case "qtyinstock":
					product.Qtyinstock = f
				case "purchase_cost":
					product.PurchaseCost = f
				}
			}
		case "reorderlevel", "qtyindemand":
			if i, ok := value.(int); ok {
				switch key {
				case "reorderlevel":
					product.Reorderlevel = i
				case "qtyindemand":
					product.Qtyindemand = i
				}
			}
		case "starred":
			product.Starred = value.(string) == "1"
		case "taxclass":
			product.Taxclass = value.(string)
		case "usageunit":
			product.Usageunit = value.(string)
		case "description":
			product.Description = value.(string)
		case "imageattachmentids":
			product.Imageattachmentids = value.(string)
		case "discontinued":
			product.Discontinued = value.(string) == "1"
		case "label":
			product.Label = value.(string)
		case "currency_id":
			product.CurrencyId = value.(string)
		}
	}

	return product, nil
}
