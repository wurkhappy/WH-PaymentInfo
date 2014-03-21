package models

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/balanced-go"
	"strconv"
	// "log"
)

type BankAccount struct {
	ID             string `json:"id"`
	CanDebit       bool   `json:"can_debit"`
	AccountNumber  string `json:"account_number,omitempty"`
	BalancedID     string `json:"balanced_id,omitempty"`
	RoutingNumber  string `json:"routing_number,omitempty"`
	VerificationID string `json:"verification_id,omitempty"`
}

type BankAccounts []*BankAccount

func (b BankAccounts) ToJSON() []byte {
	jsonString := `[`
	for i, account := range b {
		var prefix string
		if i > 0 {
			prefix = `, `
		}
		accountJSON := prefix + `{` +
			`"id":"` + account.ID + `",` +
			`"can_debit":` + strconv.FormatBool(account.CanDebit) + `,` +
			`"account_number":"` + account.AccountNumber + `",` +
			`"routing_number":"` + account.RoutingNumber + `"}`
		jsonString += accountJSON
	}
	jsonString += `]`
	return []byte(jsonString)
}

func NewBankAccount() *BankAccount {
	id, _ := uuid.NewV4()
	return &BankAccount{
		ID: id.String(),
	}
}

func (b *BankAccount) ConfirmVerification(amount1 float64, amount2 float64) error {
	balVerification := new(balanced.Verification)
	balVerification.ID = b.VerificationID
	bErrors := balVerification.Confirm(int(amount1*100), int(amount2*100))
	if len(bErrors) > 0 {
		return formatBalancedErrors(bErrors)
	}

	if balVerification.VerificationStatus == "succeeded" {
		b.CanDebit = true
	} else {
		return fmt.Errorf("%s", "Account not verified")
	}
	return nil
}

func (a *BankAccount) ConvertFromBalancedAccount(balAccount *balanced.BankAccount) {
	accountNumberLength := len(balAccount.AccountNumber)
	a.AccountNumber = balAccount.AccountNumber[accountNumberLength-4 : accountNumberLength]
	a.BalancedID = balAccount.ID
	a.RoutingNumber = balAccount.RoutingNumber
	a.CanDebit = balAccount.CanDebit
}

func (a *BankAccount) ConvertToBalancedAccount() (balAccount *balanced.BankAccount) {
	balAccount = new(balanced.BankAccount)
	balAccount.ID = a.BalancedID
	return balAccount
}
