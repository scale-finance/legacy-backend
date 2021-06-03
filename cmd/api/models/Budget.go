package models

import (
	application "github.com/elopez00/scale-backend/pkg/app"
)

type WhiteListItem struct {
	Category	string		`json:"category"` // name of category
	Name		string		`json:"name"` // name of company being listed under this category
}

type Category struct {
	Name		string		`json:"name"` // name of category
	Budget		float64		`json:"budget"` // amount of money budgeted towards this category
}

// A budget is a combination of categories
type Budget struct {
	Categories	[]Category			`json:"categories"`
	WhiteList 	[]WhiteListItem		`json:"whitelist"`
}

// This function adds categories whit their specific whitelist items to the database. If there
// is an error with the connection with the database or failure to insert, it will be reflected
// as an error return value.
func (b *Budget) Create(app *application.App, userId string) error {
	query := "INSERT INTO categories(id, name, budget) VALUES "
	vals := []interface{} {}

	for _, category := range b.Categories {
		query += " (?,?,?),"
		vals = append(vals, userId, category.Name, category.Budget)
	}

	// prepare statement
	query = query[0:len(query)-1] // trim last comma
	stmt, err := app.DB.Client.Prepare(query)
	if err != nil {
		return err
	}

	// execute query
	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}

	// add whitelist
	if err := AddWhiteList(app, userId, b.WhiteList); err != nil {
		return err
	}

	return nil
}

// Gets all the white list items and inserts them to the database. If the function fails
// due to the databse connection or query execution, an error will be returned that reflects
// this
func AddWhiteList(app *application.App, userId string, whitelist []WhiteListItem) error {
	// there might not be items that needs to be whitelisted, if this is the case return nil
	if whitelist == nil {
		return nil
	}

	query := "INSERT INTO whitelist(id, category, name) VALUES "
	vals := []interface{} {}

	for _, item := range whitelist {
		query += " (?,?,?),"
		vals = append(vals, userId, item.Category, item.Name)
	}

	// prepare statement
	query = query[0:len(query)-1]
	stmt, err := app.DB.Client.Prepare(query)
	if err != nil {
		return err
	}

	// execute query
	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}

	return nil
}