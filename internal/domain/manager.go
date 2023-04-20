package domain

type Manager struct {
	Id                string `json:"id"`
	UserName          string `json:"user_name"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email2"`
	Title             string `json:"title"`
	PhoneWork         string `json:"phone_work"`
	Department        string `json:"department"`
	Description       string `json:"description"`
	AddressStreet     string `json:"address_street"`
	AddressCity       string `json:"address_city"`
	AddressState      string `json:"address_state"`
	AddressPostalcode string `json:"address_postalcode"`
	AddressCountry    string `json:"address_country"`
	Image             string `json:"imagename"`
}

var MockedManager = Manager{
	Id:                "19x1",
	UserName:          "Administrator",
	FirstName:         "Administrator",
	LastName:          "Administrator",
	Email:             "info@itvolga.com",
	Title:             "Manager",
	PhoneWork:         "+4915211100235",
	Department:        "Management",
	Description:       "This is description for administrator user",
	AddressStreet:     "",
	AddressCity:       "",
	AddressState:      "",
	AddressPostalcode: "",
	AddressCountry:    "",
	Image:             "",
}

func ConvertMapToManager(m map[string]any) Manager {
	manager := Manager{}

	for k, v := range m {
		switch k {
		case "id":
			manager.Id = v.(string)
		case "user_name":
			manager.UserName = v.(string)
		case "first_name":
			manager.FirstName = v.(string)
		case "last_name":
			manager.LastName = v.(string)
		case "email2":
			manager.Email = v.(string)
		case "title":
			manager.Title = v.(string)
		case "phone_work":
			manager.PhoneWork = v.(string)
		case "department":
			manager.Department = v.(string)
		case "description":
			manager.Description = v.(string)
		case "address_street":
			manager.AddressStreet = v.(string)
		case "address_city":
			manager.AddressCity = v.(string)
		case "address_state":
			manager.AddressState = v.(string)
		case "address_postalcode":
			manager.AddressPostalcode = v.(string)
		case "address_country":
			manager.AddressCountry = v.(string)
		case "imagename":
			manager.Image = v.(string)
		}
	}

	return manager
}
