package domain

import (
	"encoding/json"
	"fmt"
	"github.com/semelyanov86/vtiger-portal/internal/utils"
	"strconv"
	"time"
)

type LineItem struct {
	ParentID        string       `json:"parent_id"`
	ProductID       string       `json:"productid"`
	SequenceNo      InvoiceInt   `json:"sequence_no"`
	Quantity        InvoiceFloat `json:"quantity"`
	ListPrice       InvoiceFloat `json:"listprice"`
	DiscountPercent string       `json:"discount_percent"`
	DiscountAmount  string       `json:"discount_amount"`
	Comment         string       `json:"comment"`
	Description     string       `json:"description"`
	IncrementOnDel  InvoiceBool  `json:"incrementondel"`
	Tax1            InvoiceFloat `json:"tax1"`
	Tax2            InvoiceFloat `json:"tax2"`
	Tax3            InvoiceFloat `json:"tax3"`
	Image           string       `json:"image"`
	PurchaseCost    InvoiceFloat `json:"purchase_cost"`
	Margin          InvoiceFloat `json:"margin"`
	ID              string       `json:"id"`
	ProductName     string       `json:"product_name"`
	EntityType      string       `json:"entity_type"`
	Deleted         InvoiceBool  `json:"deleted"`
}

type Invoice struct {
	Subject                string          `json:"subject"`
	SalesOrderID           string          `json:"salesorder_id"`
	CustomerNo             string          `json:"customerno"`
	ContactID              string          `json:"contact_id"`
	InvoiceDate            InvoiceDate     `json:"invoicedate"`
	DueDate                InvoiceDate     `json:"duedate"`
	VtigerPurchaseOrder    string          `json:"vtiger_purchaseorder"`
	TxtAdjustment          InvoiceFloat    `json:"txtAdjustment"`
	SalesCommission        InvoiceFloat    `json:"salescommission"`
	ExciseDuty             InvoiceFloat    `json:"exciseduty"`
	HdnSubTotal            InvoiceFloat    `json:"hdnSubTotal"`
	HdnGrandTotal          InvoiceFloat    `json:"hdnGrandTotal"`
	HdnTaxType             string          `json:"hdnTaxType"`
	HdnDiscountPercent     string          `json:"hdnDiscountPercent"`
	HdnDiscountAmount      string          `json:"hdnDiscountAmount"`
	HdnS_HAmount           InvoiceFloat    `json:"hdnS_H_Amount"`
	AccountID              string          `json:"account_id"`
	InvoiceStatus          string          `json:"invoicestatus"`
	AssignedUserID         string          `json:"assigned_user_id"`
	CreatedTime            InvoiceDateTime `json:"createdtime"`
	ModifiedTime           InvoiceDateTime `json:"modifiedtime"`
	ModifiedBy             string          `json:"modifiedby"`
	CurrencyID             string          `json:"currency_id"`
	ConversionRate         InvoiceFloat    `json:"conversion_rate"`
	BillStreet             string          `json:"bill_street"`
	ShipStreet             string          `json:"ship_street"`
	BillCity               string          `json:"bill_city"`
	ShipCity               string          `json:"ship_city"`
	BillState              string          `json:"bill_state"`
	ShipState              string          `json:"ship_state"`
	BillCode               string          `json:"bill_code"`
	ShipCode               string          `json:"ship_code"`
	BillCountry            string          `json:"bill_country"`
	ShipCountry            string          `json:"ship_country"`
	BillPOBox              string          `json:"bill_pobox"`
	ShipPOBox              string          `json:"ship_pobox"`
	Description            string          `json:"description"`
	TermsConditions        string          `json:"terms_conditions"`
	InvoiceNo              string          `json:"invoice_no"`
	PreTaxTotal            InvoiceFloat    `json:"pre_tax_total"`
	Received               InvoiceFloat    `json:"received"`
	Balance                InvoiceFloat    `json:"balance"`
	HdnS_H_Percent         InvoiceFloat    `json:"hdnS_H_Percent"`
	PotentialID            string          `json:"potential_id"`
	Source                 string          `json:"source"`
	Starred                InvoiceBool     `json:"starred"`
	Tags                   string          `json:"tags"`
	RegionID               string          `json:"region_id"`
	ID                     string          `json:"id"`
	Label                  string          `json:"label"`
	ShippingHandling       InvoiceFloat    `json:"shipping_&_handling"`
	ShippingHandlingSHTax1 InvoiceFloat    `json:"shipping_&_handling_shtax1"`
	ShippingHandlingSHTax2 InvoiceFloat    `json:"shipping_&_handling_shtax2"`
	ShippingHandlingSHTax3 InvoiceFloat    `json:"shipping_&_handling_shtax3"`
	LineItems              []LineItem      `json:"LineItems,omitempty"`
	LineItemsFinalDetails  map[string]any  `json:"LineItems_FinalDetails,omitempty"`
	Currency               Currency        `json:"currency"`
}

type InvoiceDate time.Time

type InvoiceDateTime time.Time

type InvoiceBool bool

type InvoiceFloat float64

type InvoiceInt int

func (i *InvoiceDate) UnmarshalJSON(jsonValue []byte) error {
	dateStr := string(jsonValue)

	if dateStr == `""` { // Check for empty string
		*i = InvoiceDate(time.Time{})
		return nil
	}

	date, err := utils.BytesToDate(jsonValue, "2006-01-02")

	if err != nil {
		return err
	}
	*i = InvoiceDate(date)
	return nil
}

func (i InvoiceDate) MarshalJSON() ([]byte, error) {
	date := time.Time(i)
	formatted := fmt.Sprintf("\"%s\"", date.Format("2006-01-02"))
	return []byte(formatted), nil
}

func (i InvoiceDateTime) MarshalJSON() ([]byte, error) {
	date := time.Time(i)
	formatted := fmt.Sprintf("\"%s\"", date.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (i *InvoiceDateTime) UnmarshalJSON(jsonValue []byte) error {
	date, err := utils.BytesToDate(jsonValue, "2006-01-02 15:04:05")

	if err != nil {
		return err
	}
	*i = InvoiceDateTime(date)
	return nil
}

func (i *InvoiceBool) UnmarshalJSON(jsonValue []byte) error {
	data, err := utils.BytesToBool(jsonValue)

	if err != nil {
		return err
	}
	*i = InvoiceBool(data)
	return nil
}

func (i *InvoiceFloat) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return utils.ErrInvalidRuntimeFormat
	}
	if unquotedJSONValue == "" {
		*i = InvoiceFloat(0)
		return nil
	}
	value, err := strconv.ParseFloat(unquotedJSONValue, 64)
	if err != nil {
		return err
	}
	*i = InvoiceFloat(value)
	return nil
}

func (i *InvoiceInt) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return utils.ErrInvalidRuntimeFormat
	}
	if unquotedJSONValue == "" {
		*i = InvoiceInt(0)
		return nil
	}
	value, err := strconv.Atoi(unquotedJSONValue)
	if err != nil {
		return err
	}
	*i = InvoiceInt(value)
	return nil
}

func ConvertMapToInvoice(m map[string]interface{}) (Invoice, error) {
	inv := &Invoice{}

	// Decode the invoice properties
	invJSON, err := json.Marshal(m)
	if err != nil {
		return *inv, err
	}
	if err := json.Unmarshal(invJSON, inv); err != nil {
		return *inv, err
	}

	// Decode the line items
	if lineItemsMap, ok := m["LineItems"]; ok {
		lineItemsJSON, err := json.Marshal(lineItemsMap)
		if err != nil {
			return *inv, err
		}
		if err := json.Unmarshal(lineItemsJSON, &inv.LineItems); err != nil {
			return *inv, err
		}
	}

	// Decode the line item final details
	if lineItemsFinalDetailsMap, ok := m["LineItems_FinalDetails"]; ok {
		lineItemsFinalDetailsJSON, err := json.Marshal(lineItemsFinalDetailsMap)
		if err != nil {
			return *inv, err
		}
		if err := json.Unmarshal(lineItemsFinalDetailsJSON, &inv.LineItemsFinalDetails); err != nil {
			return *inv, err
		}
	}

	return *inv, nil
}
