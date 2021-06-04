package models

import (
	application "github.com/elopez00/scale-backend/pkg/app"
)

// item for white listed companies of a specific category
type WhiteListItem struct {
	Category	string		`json:"category"` // name of category
	Name		string		`json:"name"` // name of company being listed under this category
}

// category within budget
type Category struct {
	Name		string		`json:"name"` // name of category
	Budget		float64		`json:"budget"` // amount of money budgeted towards this category
}

// Object that contains category updates
type UpdateCategory struct {
	New			Category		`json:"new"`
	Original	Category		`json:"original,omitempty"`
}

// object containing whitelist update
type UpdateWhiteList struct {
	New			WhiteListItem	`json:"new"`
	Original	WhiteListItem	`json:"original,omitempty"`
}

// Object containing both category updates and whitelist updates. Neither one or the
// other are required.
type UpdateObject struct {
	Categories	[]UpdateCategory	`json:"categories,omitempty"`
	WhiteList	[]UpdateWhiteList	`json:"whitelist,omitempty"`
}

// a struct that describes what will be updated in a budget
type UpdateRequest struct {
	Change		UpdateObject		`json:"change,omitempty"`
	Add			UpdateObject		`json:"add,omitempty"`
	Remove		UpdateObject		`json:"remove,omitempty"`	
}

// A budget is a combination of categories
type Budget struct {
	Categories	[]Category			`json:"categories"`
	WhiteList 	[]WhiteListItem		`json:"whitelist,omitempty"`
	Request		UpdateRequest		`json:"request,omitempty"`
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

// Single function that handles all updates to current budget whether it be adding, removing,
// or changing. This function will only perform at most 3 queries at a time. If there is a 
// failure inserting, deleting, or updating any of the rows it will be returned as an error.
func (b *Budget) Update(app *application.App, userId string) error {
	// check if there is add updates
	if len(b.Request.Add.Categories) != 0 || len(b.Request.Add.WhiteList) != 0 {
		if err := add(app, userId, b.Request.Add); err != nil {
			return err
		}
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

// Helper function that adds from update request
// TODO see if this can be done in one query
func add(app *application.App, userId string, add UpdateObject) error {
	// queries
	queryC := "INSERT INTO categories(id, name, budget) VALUES "
	queryW := "INSERT INTO whitelist(id, category, name) VALUES "

	// values to be added
	catVals := []interface{} {}
	listVals := []interface{} {}

	// get category values
	for _, category := range add.Categories {
		queryC += " (?,?,?),"
		catVals = append(
			catVals, 
			userId, 
			category.New.Name, 
			category.New.Budget,
		)
	}

	// get list values
	for _, item := range add.WhiteList {
		queryW += " (?,?,?),"
		listVals = append(
			listVals, 
			userId, 
			item.New.Category, 
			item.New.Name,
		)
	}

	// remove trailing comma
	queryC = queryC[0:len(queryC)-1]
	queryW = queryW[0:len(queryW)-1]

	// prepare categores statements
	stmtC, err := app.DB.Client.Prepare(queryC)
	if err != nil {
		return err
	}

	// execute categories statement
	if _, err = stmtC.Exec(catVals...); err != nil {
		return err
	}

	// prepare whitelist statemnt
	stmtW, err := app.DB.Client.Prepare(queryW)
	if err != nil {
		return err
	}

	// execute whitelist statement
	if _, err = stmtW.Exec(listVals...); err != nil {
		return err
	}

	return nil
}