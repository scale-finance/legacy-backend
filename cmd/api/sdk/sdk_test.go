package sdk_test

import (
	"github.com/elopez00/scale-backend/cmd/api/models"
)

// * test variables 

var user = models.User{
	Id: "testvalue",
	FirstName: "Stan",
	LastName: "Marsh",
	Email: "smarsh@southpark.com",
	Password: "southpark",
}

var token = models.Token{
	Value: "access-sandbox-3b6a6577-4c02-4fc3-a213-b8adf828c38f",
	Id:    "nothin",
	Institution:  "institution",
}

var publicToken = models.Token{
	Value: "public-sandbox-4d532c06-b9b5-4a18-906a-df480f320cc9",
}

var budget = models.Budget{
	Categories: []models.Category{
		{Name: "shopping", Budget: 200, WhiteList: []models.WhiteListItem{
			{Category: "shopping", Name: "Calvin Klien"},
			{Category: "shopping", Name: "Best Buy"},
			{Category: "shopping", Name: "Amazon"},
		}},
		{Name: "groceries", Budget: 250, WhiteList: []models.WhiteListItem{
			{Category: "groceries", Name: "Aldi"},
			{Category: "groceries", Name: "Walmart"},
		}},
		{Name: "rent", Budget: 800, WhiteList: []models.WhiteListItem{{Category: "rent", Name: "The Rise"}}},
	},

	Request: models.UpdateRequest{
		Update: models.UpdateObject{
			Categories: []models.Category{
				{Name: "shopping", Budget: 200, Id: "qwert"},
				{Name: "groceries", Budget: 250, Id: "asdfag"},
				{Name: "rent", Budget: 800, Id: ";lkjk"},
			},
			WhiteList: []models.WhiteListItem{
				{Category: "shopping", Name: "Calvin Klien", Id: ";lkjl"},
				{Category: "shopping", Name: "Best Buy", Id: "asdfasdf"},
				{Category: "shopping", Name: "Amazon", Id: "qwerqwer"},
				{Category: "groceries", Name: "Aldi", Id: ";sdlfkgjsd"},
				{Category: "groceries", Name: "Walmart", Id: "zxcvzxc"},
				{Category: "rent", Name: "The Rise", Id: ".,mn.,n,mn"},
			},
		},
	},
}