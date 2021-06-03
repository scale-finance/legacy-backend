package models

// import (
// 	application "github.com/elopez00/scale-backend/pkg/app"
// )

type WhiteListItem struct {
	Category	string		`json:"category"`
	Name		string		`json:"name"`
}

type Category struct {
	Name		string		`json:"name"`		// name of category
	Budget		float64		`json:"budget"`
}

type Budget struct {
	Categories	[]Category	`json:"categories"`
}

// This function adds categories whit their specific whitelist items to the database. If there
// is an error with the connection with the database or failure to insert, it will be reflected
// as an error return value.
// func AddCategories(app *application.App, userId string, categories []Category, whitelist []WhiteListItem) error {
// 	query := "INSERT INTO categories(id, name, budget) VALUES "
// 	vals := []interface{} {}

// 	for _, category := range categories {
// 		query += "(?,?,?),"
// 		vals = append(vals, userId, category.Name, category.Budget)
// 	}

// 	// trim last comma
// 	query = query[0:len(query)-1]

// 	// prepare statement
// 	stmt, err := app.DB.Client.Prepare(query)
// 	if err != nil {
// 		return err
// 	}

// 	// execute query
// 	if _, err := stmt.Exec(vals...); err != nil {
// 		return err
// 	}

// 	// add whitelist
// 	if err := AddWhiteList(app, userId, whitelist); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func AddWhiteList(app *application.App, userId string, whitelist []WhiteListItem) error {

// }