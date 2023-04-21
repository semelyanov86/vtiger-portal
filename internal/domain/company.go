package domain

type Company struct {
	OrganizationName string `json:"organizationname"`
	Address          string `json:"address"`
	City             string `json:"city"`
	State            string `json:"state"`
	Country          string `json:"country"`
	Code             string `json:"code"`
	Phone            string `json:"phone"`
	Fax              string `json:"fax"`
	Website          string `json:"website"`
	Logoname         string `json:"logoname"`
	Logo             string `json:"logo"`
	Vatid            string `json:"vatid"`
	Inn              string `json:"inn"`
	Kpp              string `json:"kpp"`
	Bankaccount      string `json:"bankaccount"`
	Bankname         string `json:"bankname"`
	Bankid           string `json:"bankid"`
	Corraccount      string `json:"corraccount"`
	Director         string `json:"director"`
	Bookkeeper       string `json:"bookkeeper"`
	Entepreneur      string `json:"entepreneur"`
	Entepreneurreg   string `json:"entepreneurreg"`
	Okpo             string `json:"okpo"`
	Id               string `json:"id"`
}

var MockedCompany = Company{
	OrganizationName: "Индивидуальный предприниматель Емельянов Сергей Петрович",
	Address:          "Бульвар Юности, д. 3, кв. 251",
	City:             "Чебоксары",
	State:            "Чувашия",
	Country:          "Россия",
	Code:             "428034",
	Phone:            "+78432023893",
	Fax:              "+78352223606",
	Website:          "https://itvolga.com",
	Logoname:         "logo.png",
	Logo:             "iVBORw0KGgoAAAANSUhEUgAAAJYAAAAnCAYAAADtl7EyAAAABmJLR0QA\\/wD\\/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH4gcNDzA7sgOMDgAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUHAAAUYElEQVR42u2caZBdxXXHf6fvfdvMm33VLBpptIBWsCTAYhUgMHYIu42xnc2FHeMklbLj5UPAlbJT\\/uLEwc5SibMVVOHEBmxjgcEYZCEhJLFJQhJaR9toGUmz72\\/evX3yoe\\/sM0gjCUcfpqtu1bx7e+\\/T\\/\\/M\\/53SP3LEnVKbTdLrIyTdMy9V0+hAEyxMZ\\/eZ85UymJ3M6jRSsEfK0c0dAR\\/v5VRSLwZy5hrIyg36YIDi4EcY2IsKH2\\/B0mpJgjQSaht0W\\/POsyUIyIVRWQBiCimBFpoRkomCs\\/cA8mulHsxkklUaMce\\/CEO3rRpI5iB+bXtVLA7FklHBcSFIFXwRjLZftepe6ww3IFFCkvaiY7ctX0puTnljw4kLfz54gu\\/Zn5D76I0xlHeIJ2Yad9H3nC6T+\\/LvErr4VtRqB2LBgqw6jnBj3Uq2CuN9qIzCcYCOoBTEO1lVHgKadRsjfCXk30boUdLbxjW9\\/jXRb65RIW5hM8b1v\\/T3vXnP9hAIpCEYEjMEABkUQPAAj0TtQFHxDeGQ\\/weG9iBj8OYvwqurQMGTgrdcg0098xSo008fA5t\\/gL11JcGAHDGSGJUyd4HlzFhMc2ImprMWUVDCw7z1MQRGxxdeA8aZV8ISCJReXdRsRYjYk3daKxqamlrzsAKmBPozIhBpUZBgtRCIhi0BJkaHvxhj6X36Gvt\\/+HPE8EMi8Ysi552FiV15Pdu2zhG3NJJZ8FG0\\/Q+9T3yev8LtkN7+M9nZjO1rQgQympAIxHolEgr5n\\/gVTNQtTUIJtPon2dhGuuJmcT37ZCdd0GqsKL95uE8ATxVzADhZ1dUyGWBK99lTxRAHByjBiegLBkb30v\\/I0XlkVufc9DBa6f\\/KP9L3wJPE5iyK9Btp2Cu1oBeMhWAq\\/+UMkDu3\\/8CjBgZ0Ufv0HmLxCbGcbvfIDBEjf9wVMcTndTz1O5u3fkvzIDcQuu3IatcaR92il9CIJFigRpz5PxFMEZbBf4wzCQcQyLh\\/CUN5BxLJtZ7B9Pdgje+l4\\/BuOGwVZV7C\\/GxFBuzvo+P5fuXKJFDJYfgTNEqx7ohdeWRV+ZQ0STxCrmc3Ajs3Y9jNR+9OCNTl5vwiC5YvjP+fP0wSPSVRhhEoAEoZ4Yei4WeisDrEWz1pihaWYnDSxmnrSd\\/4hkspF+3sJ287g5xWCKiZdQNGffYewo4XWv\\/sqRgSf0YLlRX0JB22btjPQ3ITkFxE2HcPEk8QKSvCdiE+n0RyLi4ZYyCChvjDhNJN4KUQAG2D7emj\\/0d8isTiIoJleJJ6g+\\/knCZsaKXjg86Rv\\/yQ9a39O+xPfQzw\\/sv48EnVz0SCLhgF+USkGi830gQ2HOVyQRbMDbizi+uMkOEvXsz8iaG7CdrSQWnETiflLIoSdFqZRgjVajV0cB+aHNscBpJavIlYxk7H6Nmw5RdcrzxK2NIFA3ur7SSxcTtDY4Ih4ugC\\/th6vqIz83\\/ssmslgkimkuIKSL36LxJxFYBWyQvrW+0ituBlJ5oyerBmzyL\\/zDxg4vBeTziexYJnzm03zqwlU4dBWu\\/DJMQJeZBleUB2TCadaUvMWIZctGvcpc7iBntdfQDwfT0BVSVbPQmpnDQ1PrfNlpVfc4LwJgUIiRfxj96Fh5LCySu4V1wx\\/H\\/JZBUgYkKisIVlb5+oLI241DVcTCdbFJe+eUS7E0nSCpZOzFrUMkZ4RKVk9k5rv\\/DfieUhWnXmpOj6v4JBpCKZH5BmUj0Hvv0S2qO+RuuxKEtWznHoM7RiYn0asD9VBOkiujVxIHRo5WnXyVoQR6if67XmQm47W+BxCCMZznQ3sWQRDMblpah79PmpBB8JpQZqKVXixpspDLoJVCJMytTAEG4LvO3+Ute6diHt3LozRGMzebXgHtpG9+QFI5X4wT1Kg30a1Tqu9c9Q8GqGEXgTEUgz2gozvwf5M+PgGf9OvSPzbo5jW0wgW\\/911JP79MeI\\/\\/1dMNjNk2Q2VkTF1Dnrnm4\\/h7X4bCbPj84zIN1n5D\\/w22e\\/J6j2XPl+MOgd\\/T9gek9c\\/UfnJ+hc9Q+TdbVgFlalvykEOizjEkSjIMlVryUaxv0nIuwqIWrAhxgj0dmP2vgNBFtQ6\\/rNvK5ouhKrZ0HISOX0MrV8MiRzXp+MHUWMQ46HG4GX6kFONaGUt5BZGAzHIycPQ2QqJJFo1G2IJ6GpHThxEi8rBWqSnA50xC3LyoLsDOXYAzS+GqlkQBMjRfRBm0blXOCu2vdnVK4KmC9DKOmdVNux0q1a\\/GFqbXJ8r66Cw1LV5ZA9atwDicVenGOzshe50RxggJw5DXxeamw9V9U7NNx1B2pvRmrmQm4fs346mcqF2fvTtDDpzPhpPYvZvR4sroLwGTjUiraewc5Yg8UQkGxbTsAuMQfOLkOYmtGYO5Oa7MZ0+hs66HOLJCTiWwLx6w\\/Ez54c08RhcPlPAKn0FBWy6+1Msff236BTc8C2V1TTNnYenk3uyJdoIxgjsfgeSORCLOxQKA+Tl\\/4U5i9CqWdCwE3ntF\\/Anf+0m9d31yOaXIF3gBM14mC0vI0f3onmFcOfnobAENj6PvP+m222qaMVM+MTn4MwxV39ZletfWzPMqEM\\/9hloPoH8+seQV4g++Bdw5gSy5r+cYM5d6mJNe95Btm1wQqaKLr0OVt6BrP8FeD46fzFydC+s+wWsfhAtLoPm466eh76C9nQia5+BnDxk1mUQBLD2GeTIXgbRQesXw+2fRna\\/BTs2ofd8AXLTyKs\\/hco6dPZ8ZPebsGMz3P9ltLgM+dWTsOIWtLwa9ryNvLUW+eLfQCLh5sDaqI8x9OOfQzb9CuYsQa9eDa\\/+FHLzkTmLRq3ZqPNYd13vYxUMBhGZshWtWNQqYSzB0197jKe\\/+ti5EX4NEBQVQUXOrkpFkI5W2PM2LLsJWppGs58zJ5DXn4cIHRCQ3i549zVIF8I9D8Nba+HYAWT5Kpi9EHnlJ+iOTciCZbBjI8yYDbc\\/BPu2Ia89B9s3Qlm1m7zcfLjtQdj1Jmz4JXL8IMQTrq2eLmT7G9B4YPiYDkCokMpBF16FZAdg31Zk31a45jYm8yXKyL\\/6e5H33nB8Eufs5fBu2LMVrr4Flq2C155zdc6\\/wpWxFtm5GY7shWx2jCGr8N5G56uLNv\\/QN89D3nrVjbO4AmbOG1phKSyBm+6F19e4NejthlX3wRh\\/nu\\/JaAvKE6Gjv4mG1m2ENnvOQpWM5TK3ZDkpPw1Yanvfor57w1m4m6UrVsn7+XfR7+WdneuJh0VRtbB9AyRzMPWXY19fE8X67HBAWm1kGwoGCzaL7e1CSsoxeQXYVA6qiskvQNWino\\/0tGMyvdiBAaSoDJOTQqvqsALS2YpUVGMVpLAYycuDmtnYIMBketFEHA1DZM7l6K4tTnWWVyMdLRix6Kmj2HfWIX4Mmb0Q9XwIs3hiCSNTxQioKBYX\\/zSiWCzq+7D1NSe4qVwX4TCK9nWhwQBSNQtJpdCSMnRPgOntQI2igKgi6uZlMIowaAcLFtS6v0UxYt38Dn5rPoG+sw65ZjUaCY3BQsUMd4jz0C6YOQ9TUOgCXyN2hi9jYCnm+by8\\/z841LJtaraiGK6ddT831z9EcuAkDzf8PrkDZ9er6hnW1DzOprJHzhoMEnE8S8MQGhvwPvEZTConckspRgRVRcqrMbfcC5texh7a4xy2yRS2ciZ65gR6cCe0nXFc6uh+OH0ctSGmejamsARbWAqN++HwPnTfVvDjmLr5zgIVgZNH4dBu9L3Njr8Ul2KCgCAMnEpI5WBKKtHTx6IgtUAm41CnZg5SW4\\/ufifyrTqE1iCLnDqOdrqz4dLVivT3YkQIEWhuwqx+AN3yGzQM8FBsaSVhXiFs34gYD3twN+TmYcprsacaUTGYj1yPzJyH3bnFbVwZosSY5TdhisrJbl0fbUwn4BqGmHlXou1n0IZdSG+3m1cUY0PCdzcgfhxz2ycJN76INOx0SD\\/am2MY+fhAZ6bFmfImdu6PWnoy7XgCMc06oRo8efcBj6glHbREAV9zlgdMMhcpLMWbtxS\\/fjFGBZNfjEkXYsSDghIklcYLLSaRdItufLx4ivjNd2PKa2D3Vkws7nZYbzfa2IB\\/1WpiC1bg5RURu\\/3TSFE5un4N0t5C7Ka78Ocujo5NK5JXiOzZhrS34F97B35lHcbzkaJSvFQe8Y9\\/Fv\\/q1UgyFykowRPBr5qNf8V10NWBaWvGlFZi8ovwMEhBMWIM9vkn4dD7mKJS9J31sH8HXjyJFBTjLbmG2OyFSF6RG68Fv2ImsVsfQEKLrnvOBcVXfwqvshqTTGMKS\\/DEc3NRUILkFuApmFSu+4bBWEUKyzDJtPudSiNlM9DXnkN3bMFbuAJ\\/yUcxBcVIfjHS1AiNB\\/BXrMJfcBXe0mvRvdswfb2urWit5LEjwShYinsej2\\/+Ci09x4bOLZ2bRRdwZfVq7l3wCDn9R\\/nae3XnaF0Ka6u+xbrKR8\\/NRzRo92rkhY\\/8UkO\\/R\\/49dI7GjkJWPCG7\\/gWyb\\/ya1CPfRvILXZ6RVqyYEZ54t8fDhp1k1jxBbMUqYqvucu91xHnlwQsdg\\/UM1jHkyR+RZ1BTWDsu7okIAy\\/9GFMxE\\/\\/K64Y5kepw3rFjEkbPydDYI3z6oHmJjIlR4xjJxUa2O5hnbP4xdxV8LqrnXc8rxDHoQzmnUqrjq7fhxH9PlFdDRL1hwoFzX4xzkGo48TCMW0RR0HNoa9S7kXl0kv4PbvAb7gTPH\\/9tgrwT9nVsfz6or2cbx7g8nDX\\/uJCOuSDf8vC58\\/Nziv6OUhgQ\\/8j1xBYuw6Ryh3f52fo5cy7eH30dSSSRMPhwT2ElU+NXceQVt0v8utu4C6veBcb5PJl6rFCGnKpnKdjXA9mB0dkUF8rJzZ\\/aROfkQE7u1MrEE5BIRo7k\\/4fwTksT+HHIL0SO7kdnzrtkQ0zjBMuX8++qiOBHx16mqgp9OQtmqUW2v4407h+t661CRTX29k9BNhxzh0tH73AY\\/3sUH5mAcwzllWH00BFNjOVHo3jOsKIf15fJeMxI7jKynIg72Lhvm\\/uWV4hewjeE\\/LEXH8wF9NNo9Ezxvp0BjJ6NYwlcdQssXwXHDsD6X8LVt8P8pe6o8UDoFrSnE3o63ALnpp1DVBWaT7oFKZnhFqv5hCO9pTNc6KanC2bUuetfraecdz6vyHGLzjbo73H50\\/kOHRGHns3HRwhdJDDltU5Ieruhq905NZNJyC92iNPV5tocLBfPgYIiOH3MtZtfAr2d0NECiRSUVLl2BjKwZKWbjs5W55gtrYJY\\/NITLH\\/MhvPMBSBWpEq9KerCwcN9Z23YuAnUWJxQFRPzMclUdPVasQd3YXducevrxdBsBrNkJabuMsJ9W9FDuzGf+BwSZAlf+h9kwXK80krCHZvRvVvx\\/vibaEsT9qWnMFdch1mxinDXm+i+bUgyB81mwPMxV16PmXU52tFO+PyTUFyO5OajzSfAWryH\\/hLt6MBuedn5ruJJtLcHqZ2Ld\\/Wt2CP7sJteRIorIJVGyqqQuUsIn3\\/C1X3VLYTbXkfffxupm49392cJ316HnjqKd\\/+XINNH+Ooz0NuJd9+funousaM8\\/thwr1wIdYj4lTkvNTqVhiViZsPXwbS7A7tjExKL4914N8Ri2DdfxW5dj6mahbfiFsL2FnTHFndGvqwKf+lK8HzE91E\\/hrEhimL9GOLHoeUU7NqCqZuPd\\/VqtLuTYO2z6M4tSGVdNGCDWbISb\\/5SghefQtubEWux2zZAZxvezfdiiisJ39uI3f0uWlGL+D7EEpgrrsXUL3bB3Y4Whzy+D02NaMMuFybxfTda30f9OEYtwdYNkM04x62YYffLJaUKx7zwOH\\/hF42ESjys8TE2OKdCoSTxptDscLBGh8pppg\\/6etGeLsJf\\/3iY6wRZzEAGKalA519B+OYroIq37EZMbgFoiI04TPDcf6LhoJ9Hkd5OtL8Xr6waE09CUQKbzMH2dCFB1sXfjMHzPbfAEckyNiRoOYUkc\\/ELyyEWR4sqsDZEOloglY4ogGCME4yhY4v9vdhtGzDV9diTh909yyjUiIBteB86WzAVNdimxiGL+tLjWBMFPs8TsSRCq4xfwrqZP6Sq+42zOhEG\\/DwOFt0T4c+5BrvNkDU5SHFtIuXcAJ6Pv\\/IOTLoAzWbQliZMMgd6OrENOzDltWgwgB7dj8xdgqTSEVwq8dUPoq2nGNj4AqKCpPIgnkBPn8DMuwLb1Q59vUgqF+P7BE1HwY8jqfTQ7SRFMOJhCkvRlia0owWvvJqw7bSLB+YXO5UaGTtmlFNBsI0NUFBMfNFVDJw6OnxWSoCBDHbvu3jLboTmk5FguaNKl5wqNGOO8Xp4nO9dQyPu9GcoHjsrvsR7lY+cG8ppgJnCfyRRHKoYLEYsYDH5BbD4GjLbXifY9KILAtsQ7eshWVRK0LALbTlN4rYHAOj\\/zdPYfduILf0oBiUUwc\\/LJ+zvdv44DYmXVaDzryQ4uIvMmiYIslgbkFi0Am3cR\\/j+m3g1s4mVVTrHaxRr83whsewG+tetIbvhlwTxJNrTiamdS6x+PsHe7dFtJo3GLQ6FRVAbklx+44gbQtH\\/pojO8JvyahL1CxhoccaIF12qvdQES\\/7pxOiQjiceLzQ8zcYTa7F67osdNwnumvdplldci9XwQ+207e4gPHEYU16DV1Q6bHKrEna0YFtPu8N\\/sQSmoBiTLiA4fggE\\/Nq5CJA91gAKsdo5Q\\/cE43MWY\\/t7CY8fwhSX4ZXOQMOA8MwJtKvdRVaLy\\/GKywmbT2F7OvDKqzHJHFQt4bGDaKYff\\/bliOdju9oJm09CNovk5uGVVyPxBGHbGezpE3iVtZj8IjemgX7CQ3uQnDRe1WwIsgSNB5BUGn\\/mHIKjB9DeLrzyGkxBMeGpY9iOVvy6+UgiecmpQvnn44GORwQdOiYxpbCMGKbTdHLuBplETKbvyk2nC+JYYqdnYTpd9PR\\/383AOOEEyJ0AAAAASUVORK5CYII=",
	Vatid:            "",
	Inn:              "212802719116",
	Kpp:              "0",
	Bankaccount:      "40802810602500022487",
	Bankname:         "ТОЧКА ПАО БАНКА \"ФК ОТКРЫТИЕ\"",
	Bankid:           "044525999",
	Corraccount:      "30101810845250000999",
	Director:         "Емельянов С.П.",
	Bookkeeper:       "Емельянов С.П.",
	Entepreneur:      "Емельянов Сергей Петрович",
	Entepreneurreg:   "308213034700026",
	Okpo:             "0157757201",
	Id:               "23x1",
}

func ConvertMapToCompany(m map[string]any) Company {
	company := Company{}

	for k, v := range m {
		switch k {
		case "id":
			company.Id = v.(string)
		case "organizationname":
			company.OrganizationName = v.(string)
		case "address":
			company.Address = v.(string)
		case "city":
			company.City = v.(string)
		case "state":
			company.State = v.(string)
		case "country":
			company.Country = v.(string)
		case "code":
			company.Code = v.(string)
		case "phone":
			company.Phone = v.(string)
		case "fax":
			company.Fax = v.(string)
		case "website":
			company.Website = v.(string)
		case "logoname":
			company.Logoname = v.(string)
		case "logo":
			company.Logo = v.(string)
		case "vatid":
			company.Vatid = v.(string)
		case "inn":
			company.Inn = v.(string)
		case "kpp":
			company.Kpp = v.(string)
		case "bankaccount":
			company.Bankaccount = v.(string)
		case "bankname":
			company.Bankname = v.(string)
		case "bankid":
			company.Bankid = v.(string)
		case "corraccount":
			company.Corraccount = v.(string)
		case "director":
			company.Director = v.(string)
		case "bookkeeper":
			company.Bookkeeper = v.(string)
		case "entepreneur":
			company.Entepreneur = v.(string)
		case "entepreneurreg":
			company.Entepreneurreg = v.(string)
		case "okpo":
			company.Okpo = v.(string)
		}
	}

	return company
}
