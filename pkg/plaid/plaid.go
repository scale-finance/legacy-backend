package plaid

import (
	"net/http"

	"github.com/plaid/plaid-go/plaid"
)

type Plaid struct {
	Client			*plaid.Client
	RedirectURL		string
}

// plaid.Get will return a Plaid client given 
func Get(plaidConfig map[string]string) (*Plaid, error) {
	clientOptions := plaid.ClientOptions {
		ClientID: 		plaidConfig["client"],
		Secret: 		plaidConfig["secret"],
		Environment:	plaid.Sandbox,
		HTTPClient:		&http.Client{},
	}

	client, err := plaid.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	return &Plaid { Client: client, RedirectURL: plaidConfig["redirectURL"] }, nil
}