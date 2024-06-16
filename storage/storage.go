package storage

type Priority = string

type Lead struct {
	ClientID  int    `json:"client_id"`
	LeadID    string `json:"lead_id"`
	LeadStart string `json:"lead_start"`
	LeadEnd   string `json:"lead_end"`
}

type Client struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Priority     Priority `json:"priority" binding:"oneof=HIGH MEDIUM LOW"`
	LeadCapacity int      `json:"lead_capacity"`
	Leads        []Lead   `json:"leads"`
}

type ClientRequest struct {
	Name         string   `json:"name"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Priority     Priority `json:"priority" binding:"oneof=HIGH MEDIUM LOW"`
	LeadCapacity int      `json:"lead_capacity"`
}

type AssignLeadRequest struct {
	LeadStart string `json:"lead_start"`
	LeadEnd   string `json:"lead_end"`
}

var PriorityMap = map[string]int{
	"HIGH":   3,
	"MEDIUM": 2,
	"LOW":    1,
}
