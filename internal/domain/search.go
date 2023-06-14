package domain

type Search struct {
	Label  string `json:"label"`
	Crmid  string `json:"crmid"`
	Parent string `json:"parent"`
	Module string `json:"module"`
}

func ConvertMapToSearch(m map[string]any) Search {
	search := Search{}
	for k, v := range m {
		switch k {
		case "label":
			search.Label = v.(string)
		case "question":
			search.Label = v.(string)
		case "ticket_title":
			search.Label = v.(string)
		case "projectname":
			search.Label = v.(string)
		case "id":
			search.Crmid = v.(string)
		case "parent_id":
			search.Parent = v.(string)
		case "module":
			search.Parent = v.(string)
		}
	}
	return search
}
