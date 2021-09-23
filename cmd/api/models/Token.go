package models

import (
	"github.com/elopez00/scale-backend/pkg/application"
)

// TODO use prepare statement

// Token for use of plaid public token retrieval
type Token struct {
	Value string `json:"value"`
	Id    string `json:"id"`
	Institution  string `json:"name"`
}

// Add method adds the permanent plaid token and stores into the plaid tokens table with the
// same id as the user. This function accepts two strings. The first one being the
// string describing the permanent token, and the second being a string that describes
// the item ID. Any problem given by this request will be reflected by the returned error.
func (t *Token) Add(app *application.App, userId string) error {
	query := "INSERT INTO plaidtokens(id, token, itemID, institution) VALUES(?,?,?,?)"
	stmt, err := app.DB.Client.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userId, t.Value, t.Id, t.Institution)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) Update(app *application.App, userId string) error {
	query := "UPDATE plaidtokens SET token = ? WHERE itemId = ? AND id = ?"
	_, err := app.DB.Client.Exec(query, t.Value, t.Id, userId)
	if err != nil {
		return err
	}

	return nil
}

// GetTokens returns every token associated to the user in the form of a slice of Token pointers.
// Any problem with the query or database operation will be reflected as an error and the slice
// will be returned as nil.
func GetTokens(app *application.App, userId string) ([]*Token, error) {
	query := "SELECT id, token, itemID, institution FROM plaidtokens WHERE id = ?"

	// get rows from query
	rows, err := app.DB.Client.Query(query, userId)
	if err != nil {
		return nil, err
	}

	// create slice of token pointers
	tokens := make([]*Token, 0)
	var placeholder string // ! this variable is only here because I don't know how to test without it

	// loop over all the rows and create a token for each
	for rows.Next() {
		token := new(Token)
		if err := rows.Scan(&placeholder, &token.Value, &token.Id, &token.Institution); err != nil {
			return nil, err
		}

		// each token created is to be appended to the slice
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// Get will get a token from the database and return it given the user's ID and the
// token id
func (t *Token) Get(app *application.App, userId string) error {
	// get row
	query := "SELECT id, token, itemID, institution FROM plaidtokens WHERE id = ? AND itemID = ?"
	row := app.DB.Client.QueryRow(query, userId, t.Id)

	// make token
	var placeholder string // don't know how to avoid this

	if err := row.Scan(&placeholder, &t.Value, &t.Id, &t.Institution); err != nil {
		return err
	}

	return nil
}