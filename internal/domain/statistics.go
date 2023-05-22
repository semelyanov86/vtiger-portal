package domain

type Statistics struct {
	Tickets  TicketStatistics  `json:"tickets"`
	Projects ProjectStatistics `json:"projects"`
	Tasks    TaskStatistics    `json:"tasks"`
	Invoices InvoiceStatistics `json:"invoices"`
}

type TicketStatistics struct {
	Total                int     `json:"total"`
	Open                 int     `json:"Open"`
	InProgress           int     `json:"In Progress"`
	WaitForResponse      int     `json:"Wait For Response"`
	Closed               int     `json:"Closed"`
	OpenHours            float64 `json:"Open-hours"`
	OpenDays             float64 `json:"Open-days"`
	InProgressHours      float64 `json:"In Progress-hours"`
	InProgressDays       float64 `json:"In Progress-days"`
	WaitForResponseHours float64 `json:"Wait For Response-hours"`
	WaitForResponseDays  float64 `json:"Wait For Response-days"`
	ClosedHours          float64 `json:"Closed-Hours"`
	ClosedDays           float64 `json:"Closed-Days"`
}

type ProjectStatistics struct {
	Total       int     `json:"total"`
	Open        int     `json:"open"`
	OpenHours   float64 `json:"open_hours"`
	Closed      int     `json:"closed"`
	ClosedHours float64 `json:"closed_hours"`
}

type TaskStatistics struct {
	Total           int     `json:"total"`
	Open            int     `json:"open"`
	OpenHours       float64 `json:"open-hours"`
	InProgress      int     `json:"In Progress"`
	InProgressHours float64 `json:"In Progress-hours"`
	Completed       int     `json:"Completed"`
	CompletedHours  float64 `json:"Completed-hours"`
}

type InvoiceStatistics struct {
	TotalQty int     `json:"total_qty"`
	TotalSum float64 `json:"total_sum"`
	OpenQty  int     `json:"open_qty"`
	OpenSum  float64 `json:"open_sum"`
	PaidQty  int     `json:"paid_qty"`
	PaidSum  float64 `json:"paid_sum"`
}
