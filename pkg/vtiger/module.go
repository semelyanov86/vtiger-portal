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
	Name      string `json:"name"`
	Label     string `json:"label"`
	Mandatory bool   `json:"mandatory"`
	Isunique  bool   `json:"isunique"`
	Nullable  bool   `json:"nullable"`
	Editable  bool   `json:"editable"`
	Default   string `json:"default"`
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

var MockedModule = Module{
	Label:           "Test Module",
	Name:            "TestModule",
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
			Name:      "field1",
			Label:     "Field 1",
			Mandatory: true,
			Isunique:  true,
			Nullable:  false,
			Editable:  true,
			Default:   "default1",
		},
		{
			Name:      "field2",
			Label:     "Field 2",
			Mandatory: false,
			Isunique:  false,
			Nullable:  true,
			Editable:  true,
			Default:   "default2",
		},
	},
}
