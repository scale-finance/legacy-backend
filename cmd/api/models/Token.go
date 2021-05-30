package models

import (
	"fmt"

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
func (t *Token) Add(app *application.App, userId string) error {
	query := "INSERT INTO plaidtokens(id, token, itemID) VALUES(?,?,?)"
	stmt, err := app.DB.Client.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userId, t.Value, t.Id)
	if err != nil {
		return err
	}

	return nil
}

// Method that returns every token associated to the user in the form of a slice of Token pointers.
// Any problem with the query or database operatinon will be reflected as an error and the slice
// will be returned as nil.
func GetTokens(app *application.App, userId string) ([]*Token, error) {
	query := fmt.Sprintf("SELECT id, token, itemID FROM plaidtokens WHERE id=%q", userId)
	
	// get rows from query
	rows, err := app.DB.Client.Query(query)
	if err != nil {
		return nil, err
	}

	// create slice of token pointers
	tokens := make([]*Token, 0)
	var placeholder string // ! this variable is only here because I don't know how to test without it
	
	// loop over all the rows and create a token for each
	for rows.Next() {
		token := new(Token)
		if err := rows.Scan(&placeholder, &token.Value, &token.Id); err != nil {
			return nil, err
		}

		// each token created is to be appended to the slice
		tokens = append(tokens, token)
	}

	return tokens, nil
}