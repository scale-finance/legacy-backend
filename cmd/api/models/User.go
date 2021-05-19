package models

import (
	"fmt"
	"log"

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
		log.Println("Prepare failure")
		return err
	} else {
		if _, err := stmt.Exec(u.Id, u.FirstName, u.LastName, u.Email, u.Password); err != nil {
			log.Println("Execution failure", err)
		}
	}

	return nil
}

// checks to see if user exists
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

// retrieves id and password based on credentials
func (u *User) GetCredentials(app *application.App, actualUser *User) error {
	query := fmt.Sprintf("SELECT email, password, id FROM userinfo WHERE email=%q", u.Email)
	if err := app.DB.Client.QueryRow(query).Scan(&actualUser.Email, &actualUser.Password, &actualUser.Id); err != nil {
		return err
	}

	return nil
}