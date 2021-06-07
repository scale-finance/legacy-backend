package budgeting

import (
	"fmt"
	"net/http"
	"encoding/json"

	application "github.com/elopez00/scale-backend/pkg/app"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/julienschmidt/httprouter"
)

// This handler will be in charge of all forms of edits whether it be additions, deletions,
// or changes by using a request body in the format of an UpdateRequest. This function will
// be used to create and change the user's budget in the application. If there is an error,
// with a query or database connection it will be logged and returned as a JSON error.
func Update(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// extract the budget information from the request body
		var budget models.Budget
		json.NewDecoder(r.Body).Decode(&budget)

		// gets the user id extracted from authentication cookie for later
		// use in the creation of the row containing the permanent token
		userId := fmt.Sprintf("%v", r.Context().Value(models.Key("user")))

		// updates any items in the budget.Request 
		if err := budget.Update(app, userId); err != nil {
			msg := "Failed to store budget information in database"
			models.CreateError(w, http.StatusBadGateway, msg, err)
			return
		}

		msg := "Successfully created budget"
		models.CreateResponse(w, msg, nil)
	}
}