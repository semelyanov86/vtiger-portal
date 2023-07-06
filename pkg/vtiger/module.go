package vtiger

type Module struct {
	Label           string        `json:"label"`
	Name            string        `json:"name"`
	Createable      bool          `json:"createable"`
	Updateable      bool          `json:"updateable"`
	Deleteable      bool          `json:"deleteable"`
	Retrieveable    bool          `json:"retrieveable"`
	Fields          []ModuleField `json:"fields"`
	IdPrefix        string        `json:"idPrefix"`
	IsEntity        bool          `json:"isEntity"`
	AllowDuplicates bool          `json:"allowDuplicates"`
	LabelFields     string        `json:"labelFields"`
}

type ModuleField struct {
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	Mandatory bool      `json:"mandatory"`
	Isunique  bool      `json:"isunique"`
	Nullable  bool      `json:"nullable"`
	Editable  bool      `json:"editable"`
	Default   string    `json:"default"`
	Type      FieldType `json:"type"`
}

type FieldType struct {
	Name           string           `json:"name"`
	RefersTo       []string         `json:"refersTo,omitempty"`
	Format         string           `json:"format,omitempty"`
	DefaultValue   string           `json:"defaultValue,omitempty"`
	PicklistValues []PicklistValues `json:"picklistValues,omitempty"`
}

type PicklistValues struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func (f FieldType) IsPicklistExist(value string) bool {
	if value == "" {
		return true
	}
	for _, picklistValue := range f.PicklistValues {
		if picklistValue.Value == value {
			return true
		}
	}
	return false
}

var MockedModule = Module{
	Label:           "Assets",
	Name:            "Assets",
	Createable:      true,
	Updateable:      true,
	Deleteable:      true,
	Retrieveable:    true,
	IdPrefix:        "13",
	IsEntity:        true,
	AllowDuplicates: false,
	LabelFields:     "name",
	Fields: []ModuleField{
		{
			Name:      "title",
			Label:     "Title",
			Mandatory: true,
			Isunique:  true,
			Nullable:  false,
			Editable:  true,
			Default:   "default1",
			Type:      FieldType{Name: "string"},
		},
		{
			Name:      "description",
			Label:     "Description",
			Mandatory: false,
			Isunique:  false,
			Nullable:  true,
			Editable:  true,
			Default:   "",
			Type:      FieldType{Name: "text"},
		},
		{
			Name:      "related_to",
			Label:     "Related To",
			Mandatory: false,
			Isunique:  false,
			Nullable:  true,
			Editable:  true,
			Default:   "",
			Type: FieldType{
				Name:     "reference",
				RefersTo: []string{"Contacts", "Accounts"},
			},
		},
		{
			Name:      "status",
			Label:     "Status",
			Mandatory: true,
			Isunique:  false,
			Nullable:  false,
			Editable:  true,
			Default:   "",
			Type: FieldType{
				Name: "picklist",
				PicklistValues: []PicklistValues{
					{
						Label: "New",
						Value: "New",
					},
					{
						Label: "Used",
						Value: "Used",
					},
				},
			},
		},
		{
			Name:      "hours",
			Label:     "Hours",
			Mandatory: false,
			Isunique:  false,
			Nullable:  true,
			Editable:  true,
			Default:   "",
			Type: FieldType{
				Name: "double",
			},
		},
		{
			Name:      "date",
			Label:     "Date",
			Mandatory: false,
			Isunique:  false,
			Nullable:  true,
			Editable:  true,
			Default:   "",
			Type: FieldType{
				Name:   "date",
				Format: "yyyy-mm-dd",
			},
		},
	},
}
