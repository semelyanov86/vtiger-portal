package domain

import "encoding/json"

type SalesOrder struct {
	SalesorderNo              string          `json:"salesorder_no"`
	Subject                   string          `json:"subject"`
	PotentialID               string          `json:"potential_id"`
	CustomerNo                string          `json:"customerno"`
	QuoteID                   string          `json:"quote_id"`
	VtigerPurchaseOrder       string          `json:"vtiger_purchaseorder"`
	ContactID                 string          `json:"contact_id"`
	DueDate                   string          `json:"duedate"`
	Carrier                   string          `json:"carrier"`
	Pending                   string          `json:"pending"`
	SoStatus                  string          `json:"sostatus"`
	TxtAdjustment             InvoiceFloat    `json:"txtAdjustment"`
	SalesCommission           InvoiceFloat    `json:"salescommission"`
	ExciseDuty                InvoiceFloat    `json:"exciseduty"`
	HdnGrandTotal             InvoiceFloat    `json:"hdnGrandTotal"`
	HdnSubTotal               InvoiceFloat    `json:"hdnSubTotal"`
	HdnTaxType                string          `json:"hdnTaxType"`
	HdnDiscountPercent        string          `json:"hdnDiscountPercent"`
	HdnDiscountAmount         string          `json:"hdnDiscountAmount"`
	HdnS_H_Amount             InvoiceFloat    `json:"hdnS_H_Amount"`
	AccountID                 string          `json:"account_id"`
	AssignedUserID            string          `json:"assigned_user_id"`
	CreatedTime               InvoiceDateTime `json:"createdtime"`
	ModifiedTime              InvoiceDateTime `json:"modifiedtime"`
	ModifiedBy                string          `json:"modifiedby"`
	CurrencyID                string          `json:"currency_id"`
	ConversionRate            InvoiceFloat    `json:"conversion_rate"`
	BillStreet                string          `json:"bill_street"`
	ShipStreet                string          `json:"ship_street"`
	BillCity                  string          `json:"bill_city"`
	ShipCity                  string          `json:"ship_city"`
	BillState                 string          `json:"bill_state"`
	ShipState                 string          `json:"ship_state"`
	BillCode                  string          `json:"bill_code"`
	ShipCode                  string          `json:"ship_code"`
	BillCountry               string          `json:"bill_country"`
	ShipCountry               string          `json:"ship_country"`
	BillPobox                 string          `json:"bill_pobox"`
	ShipPobox                 string          `json:"ship_pobox"`
	Description               string          `json:"description"`
	TermsConditions           string          `json:"terms_conditions"`
	PaymentDuration           string          `json:"payment_duration"`
	InvoiceStatus             string          `json:"invoicestatus"`
	FromSite                  string          `json:"fromsite"`
	PreTaxTotal               InvoiceFloat    `json:"pre_tax_total"`
	HdnS_H_Percent            string          `json:"hdnS_H_Percent"`
	SpCompany                 string          `json:"spcompany"`
	CreatedUserID             string          `json:"created_user_id"`
	Source                    string          `json:"source"`
	Starred                   string          `json:"starred"`
	RegionID                  string          `json:"region_id"`
	ID                        string          `json:"id"`
	Label                     string          `json:"label"`
	ShippingAndHandling       InvoiceFloat    `json:"shipping_&_handling"`
	ShippingAndHandlingSHTax1 InvoiceFloat    `json:"shipping_&_handling_shtax1"`
	Currency                  Currency        `json:"currency"`
	LineItems                 []LineItem      `json:"LineItems,omitempty"`
	LineItemsFinalDetails     map[string]any  `json:"LineItems_FinalDetails,omitempty"`
	Invoices                  []Invoice       `json:"invoices,omitempty"`
}

func ConvertMapToSalesOrder(m map[string]any) (SalesOrder, error) {
	inv := &SalesOrder{}

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
