package models

import (
	"fmt"
	"log"

	application "github.com/elopez00/scale-backend/pkg/app"
)

// The User struct will be used to get any information regarding the user information.
// The id in this case scenario is a primary key and will be used to retrieve other tables
// linked to the User.
type User struct {
	Id			string	`json:"id,omitempty"`
	FirstName	string	`json:"firstname,omitempty"`
	LastName	string	`json:"lastname,omitempty"`
	Email		string	`json:"email,omitempty"`
	Password	string	`json:"password,omitempty"`
}

// Method that creates user row based on current user. If there are any errors with the
// query, these issues will be returned
func (u *User) Create(app *application.App) error {
	query := "INSERT INTO userinfo(id, firstname, lastname, email, password) VALUES(?,?,?,?,?)"
	stmt, err := app.DB.Client.Prepare(query)
	if err != nil {
		log.Println("Prepare failure")
		return err
	} 
	
	if _, err := stmt.Exec(u.Id, u.FirstName, u.LastName, u.Email, u.Password); err != nil {
		log.Println("Execution failure", err)
	}

	return nil
}

// This method checks to see if current user exists in the database. Based on the
// result of this query the function will return a boolean.
// ! This function automatically assumes that errors yield false
func (u *User) Exists(app *application.App) bool {
	var test User
	query := fmt.Sprintf("SELECT firstname, email FROM userinfo WHERE email=%q", u.Email)
	if err := app.DB.Client.QueryRow(query).Scan(&test.Id, &test.Email); err != nil {
		log.Println(err)
		return false
	} else {
		log.Println(test.FirstName)
		return true
	}
}

// Method gives gets credentials found in database using current user's Email value.
// Any problem with the query or database connection will be reflected in returned error.
func (u *User) GetCredentials(app *application.App, actualUser *User) error {
	query := fmt.Sprintf("SELECT email, password, id FROM userinfo WHERE email=%q", u.Email)
	if err := app.DB.Client.QueryRow(query).Scan(&actualUser.Email, &actualUser.Password, &actualUser.Id); err != nil {
		return err
	}

	return nil
}