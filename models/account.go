package models

import (
	"github.com/wurkhappy/Balanced-go"
)

type Account struct {
	AccountNumber    string `json:"account_number,omitempty"`
	URI              string `json:"uri,omitempty"`
	RoutingNumber    string `json:"routing_number,omitempty"`
	VerificationsURI string `json:"verifications_uri,omitempty"`
	VerificationURI  string `json:"verification_uri,omitempty"`
	CreditsURI       string `json:"credits_uri,omitempty"`
}

func (a *Account) ConvertBalancedAccount(balAccount *balanced.BankAccount) {
	a.AccountNumber = balAccount.AccountNumber
	a.URI = balAccount.URI
	a.RoutingNumber = balAccount.RoutingNumber
	a.VerificationURI = balAccount.VerificationURI
	a.VerificationsURI = balAccount.VerificationsURI
	a.CreditsURI = balAccount.CreditsURI
}
