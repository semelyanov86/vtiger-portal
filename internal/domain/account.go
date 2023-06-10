package domain

import (
	"encoding/json"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

type Account struct {
	AccountName         string `json:"accountname"`
	AccountNo           string `json:"account_no"`
	Phone               string `json:"phone"`
	Website             string `json:"website"`
	Fax                 string `json:"fax"`
	TickerSymbol        string `json:"tickersymbol"`
	OtherPhone          string `json:"otherphone"`
	AccountID           string `json:"account_id"`
	Email1              string `json:"email1"`
	Employees           string `json:"employees"`
	Email2              string `json:"email2"`
	Ownership           string `json:"ownership"`
	Rating              string `json:"rating"`
	Industry            string `json:"industry"`
	SICCode             string `json:"siccode"`
	AccountType         string `json:"accounttype"`
	AnnualRevenue       string `json:"annual_revenue"`
	EmailOptOut         string `json:"emailoptout"`
	NotifyOwner         string `json:"notify_owner"`
	AssignedUserID      string `json:"assigned_user_id"`
	CreatedTime         string `json:"createdtime"`
	ModifiedTime        string `json:"modifiedtime"`
	ModifiedBy          string `json:"modifiedby"`
	BillStreet          string `json:"bill_street"`
	ShipStreet          string `json:"ship_street"`
	BillCity            string `json:"bill_city"`
	ShipCity            string `json:"ship_city"`
	BillState           string `json:"bill_state"`
	ShipState           string `json:"ship_state"`
	BillCode            string `json:"bill_code"`
	ShipCode            string `json:"ship_code"`
	BillCountry         string `json:"bill_country"`
	ShipCountry         string `json:"ship_country"`
	BillPOBox           string `json:"bill_pobox"`
	ShipPOBox           string `json:"ship_pobox"`
	Description         string `json:"description"`
	IsConvertedFromLead string `json:"isconvertedfromlead"`
	Source              string `json:"source"`
	Starred             string `json:"starred"`
	Tags                string `json:"tags"`
	ImageName           string `json:"imagename"`
	ID                  string `json:"id"`
	ImageAttachmentIDs  string `json:"imageattachmentids"`
	Label               string `json:"label"`
}

func ConvertMapToAccount(data map[string]interface{}) (Account, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return Account{}, e.Wrap("error marshalling data", err)
	}

	var account Account
	err = json.Unmarshal(jsonData, &account)
	if err != nil {
		return Account{}, e.Wrap("error unmarshal data", err)
	}

	return account, nil
}
