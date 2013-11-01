package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
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

func (a *BankAccount) ConvertBalancedAccount(balAccount *balanced.BankAccount) {
	a.AccountNumber = balAccount.AccountNumber
	a.URI = balAccount.URI
	a.RoutingNumber = balAccount.RoutingNumber
	a.VerificationURI = balAccount.VerificationURI
	a.VerificationsURI = balAccount.VerificationsURI
	a.CreditsURI = balAccount.CreditsURI
	a.CanDebit = balAccount.CanDebit
}

func DeleteBankAccount(userID string, accountID string, ctx *DB.Context) {
	m := make(map[string]interface{})

	change := mgo.Change{
		Update:    bson.M{"$pull": bson.M{"accounts": bson.M{"id": accountID}}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("usersbal")
	query := coll.Find(bson.M{
		"_id": userID,
	})

	user := new(User)
	query.One(&user)
	balBankAccount := new(balanced.BankAccount)

	for _, bankAccount := range user.Accounts {
		if bankAccount.ID == accountID {
			balBankAccount.URI = bankAccount.URI
			bError := balBankAccount.Delete()
			if bError != nil {
				return
			}

			_, _ = query.Apply(change, &m)
		}
	}
}
