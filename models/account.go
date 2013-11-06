package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/Balanced-go"
	"fmt"
	// "log"
)

type BankAccount struct {
	ID               string `json:"id"`
	CanDebit         bool   `json:"can_debit"`
	AccountNumber    string `json:"account_number,omitempty"`
	URI              string `json:"uri,omitempty"`
	RoutingNumber    string `json:"routing_number,omitempty"`
	VerificationsURI string `json:"verifications_uri,omitempty"`
	VerificationURI  string `json:"verification_uri,omitempty"`
	CreditsURI       string `json:"credits_uri,omitempty"`
}

func NewBankAccount() *BankAccount {
	id, _ := uuid.NewV4()
	return &BankAccount{
		ID: id.String(),
	}
}

func (b *BankAccount) ConfirmVerification(amount1 float64, amount2 float64) error {
	balAccount := new(balanced.BankAccount)
	balAccount.VerificationURI = b.VerificationURI
	verification, bError := balAccount.ConfirmVerification(amount1, amount2)
	if bError != nil {
		return fmt.Errorf("%s", bError.Description)
	}

	if verification.State == "verified" {
		b.CanDebit = true
	} else {
		return fmt.Errorf("%s", "Account not verified")
	}
	return nil
}

func (a *BankAccount) ConvertFromBalancedAccount(balAccount *balanced.BankAccount) {
	a.AccountNumber = balAccount.AccountNumber
	a.URI = balAccount.URI
	a.RoutingNumber = balAccount.RoutingNumber
	a.VerificationURI = balAccount.VerificationURI
	a.VerificationsURI = balAccount.VerificationsURI
	a.CreditsURI = balAccount.CreditsURI
	a.CanDebit = balAccount.CanDebit
}

func (a *BankAccount) ConvertToBalancedAccount() (balAccount *balanced.BankAccount) {
	balAccount = new(balanced.BankAccount)
	balAccount.AccountNumber = a.AccountNumber
	balAccount.URI = a.URI
	balAccount.RoutingNumber = a.RoutingNumber
	balAccount.VerificationURI = a.VerificationURI
	balAccount.VerificationsURI = a.VerificationsURI
	balAccount.CreditsURI = a.CreditsURI
	balAccount.CanDebit = a.CanDebit
	return balAccount
}
