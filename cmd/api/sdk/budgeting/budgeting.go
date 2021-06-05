package budgeting

import (
	"fmt"
	"net/http"
	"encoding/json"

	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/julienschmidt/httprouter"
)

// This function initially creates a budget for the user and adds their categories 
// and respective whitelists to the database. If there is a failure in the insertion
// or execution of database queries will result in a JSON error.
func Create(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// extract the budget information from the request body
		var budget models.Budget
		json.NewDecoder(r.Body).Decode(&budget)

		// gets the user id extracted from authentication cookie for later
		// use in the creation of the row containing the permanent token
		userId := fmt.Sprintf("%v", r.Context().Value(models.Key("user")))

		// add categories to database 
		if err := budget.Update(app, userId); err != nil {
			msg := "Failed to store budget information in database"
			models.CreateError(w, http.StatusBadGateway, msg, err)
			return
		}

		msg := "Successfully created budget"
		models.CreateResponse(w, msg, nil)
	}
}