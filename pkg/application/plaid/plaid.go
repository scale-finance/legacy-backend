package plaid

import (
	"github.com/elopez00/scale-backend/pkg/application/config"
	"net/http"

	"github.com/plaid/plaid-go/plaid"
)

type Plaid struct {
	// Client is the object that contains all database functionalities
	Client			*plaid.Client

	// RedirectURL is necessary for something, I'll figure it out
	RedirectURL		string
}

// Get will return a Plaid client given application config
func Get(config config.Config) (*Plaid, error) {
	plaidConfig := config.GetPlaid()

	// if the client id is test, we return a nil Plaid client with a valid redirect uri
	// for test purposes
	if plaidConfig["client"] == "test" {
		return &Plaid { Client: nil, RedirectURL: plaidConfig["redirectUri"] }, nil
	}

	// else we establish plaid options
	clientOptions := plaid.ClientOptions {
		ClientID: 		plaidConfig["client"],
		Secret: 		plaidConfig["secret"],
		Environment:	plaid.Sandbox,
		HTTPClient:		&http.Client{},
	}

	// and instantiate a new client
	client, err := plaid.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	return &Plaid { Client: client, RedirectURL: plaidConfig["redirectURL"] }, nil
}