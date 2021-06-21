package models

import (
	"github.com/plaid/plaid-go/plaid"
)

// define a struct for balance types
type BType struct {
	Current     float64 `json:"current"`
	Name        string  `json:"name"`
	Limit       float64 `json:"limit,omitempty"`
	Institution string  `json:"institution"`
	Mask        string  `json:"mask"`
	Id          string  `json:"id"`
	Due         float64 `json:"due,omitempty"`
	PaymentDate string  `json:"paymentDate,omitempty"`
}

// define a struct for balance totals
type BTotal struct {
	Liquid float64 `json:"liquid"`
	Credit float64 `json:"credit"`
	Loan   float64 `json:"loan"`
	Total  float64 `json:"total"`
}

// define a struct for the balance response
type Balance struct {
	Liquid []BType `json:"liquid"`
	Credit []BType `json:"credit"`
	Loan   []BType `json:"loan"`
	Net    BTotal  `json:"net"`
}

type PlaidLiabilities struct {
	Student  []plaid.StudentLoanLiability `json:"student"`
	Credit   []plaid.CreditLiability      `json:"credit"`
	Mortgage []plaid.MortgageLiability    `json:"mortgage"`
}

// Given the institution, account, and liability, this function will add that balance to the
// balance object.
func (b *Balance) AddBalance(institution string, account plaid.Account, liabilities *PlaidLiabilities) {
	switch account.Type {
	case "depository": 
		{
			b.Liquid = append(b.Liquid, BType{
				Institution: institution,
				Current:     account.Balances.Current,
				Id:          account.AccountID,
				Name:        account.Name,
				Mask:        account.Mask,
				Limit:       account.Balances.Limit,
			})
			b.Net.Liquid += account.Balances.Current
			b.Net.Total += account.Balances.Current
		}
	case "credit":
		{
			if len(liabilities.Credit) != 0 {
				b.Credit = append(b.Credit, BType{
					Institution: institution,
					Current:     account.Balances.Current,
					Id:          account.AccountID,
					Name:        account.Name,
					Mask:        account.Mask,
					Limit:       account.Balances.Limit,
					Due:         liabilities.Credit[0].LastPaymentAmount,
					PaymentDate: liabilities.Credit[0].LastPaymentDate,
				})
		
				b.Net.Credit += account.Balances.Current
				b.Net.Total -= account.Balances.Current
		
				// push back liabilities
				liabilities.Credit = liabilities.Credit[1:len(liabilities.Credit)]
			}
		}
	default:
		{
			if account.Subtype == "student" {
				if len(liabilities.Student) != 0 {
					b.Loan = append(b.Loan, BType{
						Institution: institution,
						Current:     account.Balances.Current,
						Id:          account.AccountID,
						Name:        account.Name,
						Mask:        account.Mask,
						Limit:       account.Balances.Limit,
						Due:         liabilities.Student[0].LastPaymentAmount,
						PaymentDate: liabilities.Student[0].LastPaymentDate,
					})
			
					b.Net.Loan += account.Balances.Current
					b.Net.Total -= account.Balances.Current
			
					// push back liabilities
					liabilities.Student = liabilities.Student[1:len(liabilities.Student)]
				}
			} else {
				if len(liabilities.Mortgage) != 0 {
					b.Loan = append(b.Loan, BType{
						Institution: institution,
						Current:     account.Balances.Current,
						Id:          account.AccountID,
						Name:        account.Name,
						Mask:        account.Mask,
						Limit:       account.Balances.Limit,
						Due:         liabilities.Mortgage[0].LastPaymentAmount,
						PaymentDate: liabilities.Mortgage[0].LastPaymentDate,
					})
			
					b.Net.Loan += account.Balances.Current
					b.Net.Total -= account.Balances.Current
			
					// push back liabilities
					liabilities.Mortgage = liabilities.Mortgage[1:len(liabilities.Mortgage)]
				}
			}
		}
	}
}