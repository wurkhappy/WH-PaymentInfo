package models

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"log"
)

type User struct {
	ID        string       `json:"id" bson:"_id"`
	URI       string       `json:"uri"`
	DebitsURI string       `json:"debitsURI"`
	Cards     Cards        `json:"cards"`
	Accounts  BankAccounts `json:"accounts"`
}

func (u *User) GetBankAccount(accountID string) *BankAccount {
	var bankAccount *BankAccount
	for _, account := range u.Accounts {
		if account.ID == accountID {
			bankAccount = account
		}
	}
	return bankAccount
}

func (u *User) GetCard(cardID string) *Card {
	var card *Card
	for _, c := range u.Cards {
		if c.ID == cardID {
			card = c
		}
	}
	return card
}

func CreateUserWithID(id string) (u *User, err error) {
	u = new(User)
	u.ID = id

	userBal := new(balanced.Customer)
	bError := userBal.Create()
	if bError != nil {
		return nil, fmt.Errorf("%s", bError.Description)
	}
	u.ConvertFromBalanced(userBal)
	return u, nil
}

func (u *User) Save() (err error) {
	jsonByte, _ := json.Marshal(u)
	_, err = DB.UpsertUser.Query(u.ID, string(jsonByte))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (u *User) ConvertFromBalanced(bal *balanced.Customer) {
	u.URI = bal.URI
	u.DebitsURI = bal.DebitsURI
}

func FindUserByID(id string) (u *User, err error) {
	var s string
	err = DB.FindUserByID.QueryRow(id).Scan(&s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &u)
	return u, nil
}

func (u *User) AddCreditCard(card *Card) error {
	go func(user *User, c *Card) {
		userBal := new(balanced.Customer)
		userBal.URI = user.URI
		balCard := c.ConvertToBalancedCard()
		bError := userBal.AddCreditCard(balCard)
		if bError != nil {
			fmt.Printf("add cc err %s user id: %s", bError, user.ID)
		}
	}(u, card)
	u.Cards = append(u.Cards, card)
	return nil
}

func (u *User) DeleteCard(cardID string) error {
	for i, card := range u.Cards {
		if card.ID == cardID {
			balCard := new(balanced.Card)
			balCard.URI = card.URI
			bError := balCard.Delete()
			if bError != nil {
				return fmt.Errorf("%s", bError.Description)
			}
			u.Cards = append(u.Cards[:i], u.Cards[i+1:]...)
			err := u.Save()
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func (u *User) AddBankAccount(ba *BankAccount) error {
	userBal := new(balanced.Customer)
	userBal.URI = u.URI
	balAccount := ba.ConvertToBalancedAccount()
	bError := userBal.AddBankAccount(balAccount)
	if bError != nil {
		return fmt.Errorf("%s", bError.Description)
	}
	_, bError = balAccount.Verify()
	if bError != nil {
		return fmt.Errorf("%s", bError.Description)
	}
	ba.VerificationURI = balAccount.VerificationURI
	u.Accounts = append(u.Accounts, ba)
	return nil
}

func (u *User) DeleteBankAccount(accountID string) error {
	for i, account := range u.Accounts {
		if account.ID == accountID {
			balAccount := new(balanced.BankAccount)
			balAccount.URI = account.URI
			bError := balAccount.Delete()
			if bError != nil {
				return fmt.Errorf("%s", bError.Description)
			}
			u.Accounts = append(u.Accounts[:i], u.Accounts[i+1:]...)
			err := u.Save()
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
