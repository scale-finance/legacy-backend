package models_test

import (
	"testing"

	// "github.com/elopez00/scale-backend/pkg/test"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/plaid/plaid-go/plaid"
)

var testAccount = plaid.GetLiabilitiesResponse {
	Liabilities: models.PlaidLiabilities {
		Credit: []plaid.CreditLiability{
			{
				AccountID: "dVzbVMLjrxTnLjX4G66XUp5GLklm4oiZy88yK",
				APRs: []plaid.APR {
					{
						APRPercentage: 15.24,
						APRType: "balance_transfer_apr",
						BalanceSubjectToAPR: 1562.32,
						InterestChargeAmount: 130.22,
				  	},
				  	{
						APRPercentage: 27.95,
						APRType: "cash_apr",
						BalanceSubjectToAPR: 56.22,
						InterestChargeAmount: 14.81,
					},
					{
						APRPercentage: 12.5,
						APRType: "purchase_apr",
						BalanceSubjectToAPR: 157.01,
						InterestChargeAmount: 25.66,
			    	},
				  	{
						APRPercentage: 0,
						APRType: "special",
						BalanceSubjectToAPR: 1000,
						InterestChargeAmount: 0,
				  	},
				},
				IsOverdue: false,
				LastPaymentAmount: 168.25,
				LastPaymentDate: "2019-05-22",
				LastStatementIssueDate: "2019-05-28",
				MinimumPaymentAmount: 20,
				NextPaymentDueDate: "2020-05-28",
			},
		},
		Student: []plaid.StudentLoanLiability{
			{
				AccountID: "Pp1Vpkl9w8sajvK6oEEKtr7vZxBnGpf7LxxLE",
        		AccountNumber: "4277075694",
        		DisbursementDates: []string { "2002-08-28" },
        		ExpectedPayoffDate: "2032-07-28",
        		Guarantor: "DEPT OF ED",
        		InterestRatePercentage: 5.25,
        		IsOverdue: false,
        		LastPaymentAmount: 138.05,
        		LastPaymentDate: "2019-04-22",
        		LastStatementIssueDate: "2019-04-28",
        		LoanName: "Consolidation",
        		LoanStatus: plaid.StudentLoanStatus{
        			EndDate: "2032-07-28",
        			Type: "repayment",
        		},
        		MinimumPaymentAmount: 25,
        		NextPaymentDueDate: "2019-05-28",
        		OriginationDate: "2002-08-28",
        		OriginationPrincipalAmount: 25000,
        		OutstandingInterestAmount: 6227.36,
        		PaymentReferenceNumber: "4277075694",
        		PSLFStatus: plaid.PSLFStatus{
        			EstimatedEligibilityDate: "2021-01-01",
        			PaymentsMade: 200,
        			PaymentsRemaining: 160,
        		},
        		RepaymentPlan: plaid.StudentLoanRepaymentPlan{
        			Description: "Standard Repayment",
        			Type: "standard",
        		},
        		SequenceNumber: "1",
        		ServicerAddress: plaid.StudentLoanServicerAddress{
        			City: "San Matias",
        			Country: "US",
        			PostalCode: "99415",
        			Region: "CA",
        			Street: "123 Relaxation Road",
        		},
        		YTDInterestPaid: 280.55,
        		YTDPrincipalPaid: 271.65,
			},
		},
		Mortgage: []plaid.MortgageLiability{
			{
				AccountID: "BxBXxLj1m4HMXBm9WZJyUg9XLd4rKEhw8Pb1J",
				AccountNumber: "3120194154",
				CurrentLateFee: 25,
				EscrowBalance: 3141.54,
				HasPmi: true,
				HasPrepaymentPenalty: true,
				InterestRate: plaid.MortgageInterestRate{
					Percentage: 3.99,
					Type: "fixed",
				},
				LastPaymentAmount: 3141.54,
				LastPaymentDate: "2019-08-01",
				LoanTerm: "30 year",
				LoanTypeDescription: "conventional",
				MaturityDate: "2045-07-31",
				NextMonthlyPayment: 3141.54,
				NextPaymentDueDate: "2019-11-15",
				OriginationDate: "2015-08-01",
				OriginationPrincipalAmount: 425000,
				PastDueAmount: 2304,
				PropertyAddress: plaid.MortgagePropertyAddress{
					City: "Malakoff",
					Country: "US",
					PostalCode: "14236",
					Region: "NY",
					Street: "2992 Cameron Road",
				},
				YtdInterestPaid: 12300.4,
				YtdPrincipalPaid: 12340.5,
			},
		},
	},

	Accounts: []plaid.Account{
		{
			AccountID: "dVzbVMLjrxTnLjX4G66XUp5GLklm4oiZy88yK",
			Balances: plaid.AccountBalances {
				Current: 410,
				ISOCurrencyCode: "USD",
				Limit: 2000,
			},
			Mask: "3333",
			Name: "Plaid Credit Card",
			OfficialName: "Plaid Diamond 12.5% APR Interest Credit Card",
			Subtype: "credit card",
			Type: "credit",
		},
		{
			AccountID: "Pp1Vpkl9w8sajvK6oEEKtr7vZxBnGpf7LxxLE",
			Balances: plaid.AccountBalances{
			  Current: 65262,
			  ISOCurrencyCode: "USD",
			},
			Mask: "7777",
			Name: "Plaid Student Loan",
			Subtype: "student",
			Type: "loan",
		},
		{
			AccountID: "BxBXxLj1m4HMXBm9WZJyUg9XLd4rKEhw8Pb1J",
			Balances: plaid.AccountBalances{
			  Current: 56302.06,
			  ISOCurrencyCode: "USD",
			},
			Mask: "8888",
			Name: "Plaid Mortgage",
			Subtype: "mortgage",
			Type: "loan",
		},
		{
			AccountID: "dVzbVMLjrxTnLjX4G66XUp5GLklm4oiZy88yK",
			Balances: plaid.AccountBalances {
				Current: 410,
				ISOCurrencyCode: "USD",
				Limit: 2000,
			},
			Mask: "3333",
			Name: "Plaid Debit Card",
			OfficialName: "Plaid Diamond 12.5% APR Interest Credit Card",
			Subtype: "debit card",
			Type: "depository",
		},
	},
}

var testBalance = models.Balance {
	Credit: []models.BType{
		{
			Current: 410,
			Name: "Plaid Credit Card",
			Limit: 2000,
			Institution: "Bank 2",
			Mask: "3333",
			Id: "PVK7RNneVGtZdnNVQavJIPxRnQoMQdC7rXVXB",
		},
	},
}

func getCopies() (models.Balance, []plaid.Account, models.PlaidLiabilities) {
	var newBal models.Balance
	newBal.Credit = make([]models.BType, len(testBalance.Credit))
	copy(newBal.Credit, testBalance.Credit)

	var newLiabilities models.PlaidLiabilities
	newLiabilities.Credit = make([]plaid.CreditLiability, len(testAccount.Liabilities.Credit))
	newLiabilities.Student = make([]plaid.StudentLoanLiability, len(testAccount.Liabilities.Student))
	newLiabilities.Mortgage = make([]plaid.MortgageLiability, len(testAccount.Liabilities.Mortgage))
	copy(newLiabilities.Credit, testAccount.Liabilities.Credit)
	copy(newLiabilities.Student, testAccount.Liabilities.Student)
	copy(newLiabilities.Mortgage, testAccount.Liabilities.Mortgage)

	newAccounts := make([]plaid.Account, len(testAccount.Accounts))
	copy(newAccounts, testAccount.Accounts)

	return newBal, newAccounts, newLiabilities
}

func TestAddBalanceCredit(t *testing.T) {
	balance, accounts, liabilites := getCopies()

	balance.AddBalance("Bank 1", accounts[0], &liabilites)

	if len(balance.Credit) != 2 {
		t.Fatal("No credit was added")
	}

	if len(liabilites.Credit) != 0 {
		t.Fatal("This was supposed to remove the first element")
	}

	if balance.Credit[0].Current == 0 {
		t.Fatal("Credit account was successfully added, however, the account added was empty")
	}
}

func TestAddBalanceStudentLoan(t *testing.T) {
	balance, accounts, liabilities := getCopies()

	balance.AddBalance("Bank 2", accounts[1], &liabilities)

	if len(balance.Loan) != 1 {
		t.Fatal("No student loan added")
	}

	if len(liabilities.Student) != 0 {
		t.Fatal("This was supposed to remove the first element")
	}

	if balance.Loan[0].Current == 0{
		t.Fatal("Student loan account was successfully added, however, the account was empty")
	}
}

func TestAddBalanceMortgageLoan(t *testing.T) {
	balance, accounts, liabilities := getCopies()

	balance.AddBalance("Bank 3", accounts[2], &liabilities)

	if len(balance.Loan) != 1 {
		t.Fatal("No Mortgage loan added")
	}

	if len(liabilities.Mortgage) != 0 {
		t.Fatal("This was supposed to remove the first element")
	}

	if balance.Loan[0].Current == 0{
		t.Fatal("Mortgage loan account was successfully added, however, the account was empty")
	}
}

func TestAddAllLiabilities(t *testing.T) {
	balance, accounts, liabilities := getCopies()

	for _, account := range accounts {
		balance.AddBalance("Bank 1", account, &liabilities)
	}

	if len(liabilities.Credit) != 0 && len(liabilities.Mortgage) != 0 && len(liabilities.Credit) != 0 {
		t.Fatal("Liabilities werent removed")
	}

	if len(balance.Loan) != 2 && len(balance.Credit) != 2 {
		t.Fatal("Not all balances and liabilities were added")
	}
}