package domain

type Currency struct {
	Id             string  `json:"id"`
	CurrencyName   string  `json:"currency_name"`
	CurrencyCode   string  `json:"currency_code"`
	CurrencySymbol string  `json:"currency_symbol"`
	ConversionRate float64 `json:"conversion_rate"`
	CurrencyStatus string  `json:"currency_status"`
	Defaultid      string  `json:"defaultid"`
	Deleted        bool    `json:"deleted"`
}

var MockedCurrency = Currency{
	Id:             "21x1",
	CurrencyName:   "Euro",
	CurrencyCode:   "EUR",
	CurrencySymbol: "€",
	ConversionRate: 1.00,
	CurrencyStatus: "Active",
	Defaultid:      "-11",
	Deleted:        false,
}

func ConvertMapToCurrency(inputMap map[string]any) (Currency, error) {
	var currency Currency

	for key, value := range inputMap {
		switch key {
		case "id":
			currency.Id = value.(string)
		case "currency_name":
			currency.CurrencyName = value.(string)
		case "currency_code":
			currency.CurrencyCode = value.(string)
		case "currency_symbol":
			currency.CurrencySymbol = value.(string)
		case "currency_status":
			currency.CurrencyStatus = value.(string)
		case "defaultid":
			currency.Defaultid = value.(string)
		case "conversion_rate":
			if f, ok := value.(float64); ok {
				switch key {
				case "conversion_rate":
					currency.ConversionRate = f
				}
			}
		case "deleted":
			currency.Deleted = value.(string) == "1"
		}
	}

	return currency, nil
}
