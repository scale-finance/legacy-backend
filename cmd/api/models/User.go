package models

import (
	"fmt"

	application "github.com/elopez00/scale-backend/pkg/app"
)

type User struct {
	Id			string	`json:"id,omitempty"`
	FirstName	string	`json:"firstname,omitempty"`
	LastName	string	`json:"lastname,omitempty"`
	Email		string	`json:"email,omitempty"`
	Password	string	`json:"password,omitempty"`
}

// Creates user in database
func (u *User) Create(app *application.App) error {
	query := "INSERT INTO userinfo(id, firstname, lastname, email, password) VALUES(?,?,?,?,?)"
	if stmt, err := app.DB.Client.Prepare(query); err != nil {
		return err
	} else {
		if _, err := stmt.Exec(u.Id, u.FirstName, u.LastName, u.Email, u.Password); err != nil {
			fmt.Println("User Onboarded")
		}
	}

	return nil
}

// checks to see if user exists
func (u *User) Exists(app *application.App) bool {
	query := fmt.Sprintf("SELECT * FROM userinfo WHERE email=%q", u.Email)
	if err := app.DB.Client.QueryRow(query); err != nil {
		return false
	} else {
		return true
	}
}

// retrieves id and password based on credentials
func (u *User) GetCredentials(app *application.App, actualUser *User) error {
	query := fmt.Sprintf("SELECT password, id FROM userinfo WHERE email=%q", u.Email)
	if err := app.DB.Client.QueryRow(query).Scan(&actualUser.Password, &actualUser.Id); err != nil {
		return err
	}

	return nil
}