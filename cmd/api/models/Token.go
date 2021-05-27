package models

import (
	application "github.com/elopez00/scale-backend/pkg/app"
)

// for use of plaid public token retrieval
type Token struct {
	Value		string	`json:"value"`
	Id			string	`json:"id"`
}

// Method adds the permanent plaid token and stores into the plaidtokens table with the
// same id as the user. This function accepts two strings. The first one being the 
// string describing the permanent token, and the second being a string that describes
// the item ID. Any problem given by this request will be reflected by the returned error.
func (t *Token) Add(app *application.App, userID string) error {
	query := "INSERT INTO plaidtokens(id, token, itemID) VALUES(?,?,?)"
	stmt, err := app.DB.Client.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userID, t.Value, t.Id)
	if err != nil {
		return err
	}

	return nil
}