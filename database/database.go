package database

import (
	"encoding/json"
	"net/http"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/elopez00/scale-backend/models"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	_ "github.com/go-sql-driver/mysql"
)

// database object
var db *sql.DB

// Connects to database
func Connect(router *mux.Router) *sql.DB {
	// check environemnt
	if os.Getenv("LIVE") != "true" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal(err.Error())
		}
	}
	
	var (
		// credentials
		username string = os.Getenv("DB_USERNAME")
		password string = os.Getenv("DB_PASSWORD")
		database string = os.Getenv("DB_DATABASE")
		accesspt string = os.Getenv("DB_ACCESSPT")

		// error
		err error
	)

	signIn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, accesspt, database)
	if db, err = sql.Open("mysql", signIn); err != nil {
		fmt.Println("Error here")
		log.Fatal(err.Error())
	}

	router.HandleFunc("/onboard", handleOnboard).Methods("POST")
	router.HandleFunc("/login", handleLogin).Methods("GET")
	
	return db
}

// handles user signup and inserts
func handleOnboard(w http.ResponseWriter, r *http.Request) {
	var newUser models.User // define new user

	// decodes the JSON object and puts it into user struct
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create unique id for user
	newUser.Id = uuid.New().String()

	// prepare query
	query := "INSERT INTO userinfo(id, firstname, lastname, email, password) VALUES(?, ?, ?, ?, ?)"
	if stmt, err := db.Prepare(query); err != nil {
		fmt.Println("Hello")
		panic(err.Error())
	} else {
		if _, err = stmt.Exec(newUser.Id, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password); err != nil {
			panic(err.Error())
		} else {
			fmt.Fprintf(w, "User onboarded")
		}
	}
}

// authenticates and returns user info
func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authUser models.User 
	if err := json.NewDecoder(r.Body).Decode(&authUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(
		"SELECT id, email, firstname, lastname FROM userinfo WHERE `email`=%q AND `password`=%q",
		authUser.Email, authUser.Password,
	)	
	if result, err := db.Query(query); err != nil {
		fmt.Fprint(w, "[ERROR]: ", err.Error())
		panic(err.Error())
	} else {
		var userInfo models.User
		defer result.Close()
		
		for result.Next() {
			if err := result.Scan(&userInfo.Id, &userInfo.Email, &userInfo.FirstName, &userInfo.LastName); err != nil {
				fmt.Fprint(w, "[ERROR]: ", err.Error())
				panic(err.Error())
			} else {
				json.NewEncoder(w).Encode(userInfo)
			}
		}
	}
}