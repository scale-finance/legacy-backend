package models

type Response struct {
	Status		int		`json:"status"`
	Message		string	`json:"message"`
	Type 		string 	`json:"type"`
	Result		string	`json:"result,omitempty"`
}

type Key string

type Tkn struct {
	Public		string
}