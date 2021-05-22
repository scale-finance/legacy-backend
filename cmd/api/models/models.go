package models

// for use of json responses
type Response struct {
	Status		int		`json:"status"`
	Message		string	`json:"message"`
	Type 		string 	`json:"type"`
	Result		string	`json:"result,omitempty"`
}

// for use of context keys
type Key string

// for use of plaid public token retrieval
type Tkn struct {
	Public		string
}