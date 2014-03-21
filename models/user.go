package models

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"github.com/wurkhappy/balanced-go"
	"log"
)

type User struct {
	ID         string       `json:"id" bson:"_id"`
	BalancedID string       `json:"balanced_id"`
	IsVerified bool         `json:"isVerified"`
	Cards      Cards        `json:"cards"`
	Accounts   BankAccounts `json:"accounts"`
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
	bErrors := balanced.Create(userBal)
	if len(bErrors) > 0 {
		return nil, formatBalancedErrors(bErrors)
	}
	u.ConvertFromBalanced(userBal)
	return u, nil
}

func (u *User) UpdateWithMap(m map[string]interface{}) error {
	bUser := u.ConvertToBalanced(m)
	bErrors := balanced.Update(bUser)
	if len(bErrors) > 0 {
		return formatBalancedErrors(bErrors)
	}
	u.ConvertFromBalanced(bUser)
	return nil
}

func (u *User) Save() (err error) {
	jsonByte, _ := json.Marshal(u)
	r, err := DB.UpsertUser.Query(u.ID, string(jsonByte))
	r.Close()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (u *User) ConvertFromBalanced(bal *balanced.Customer) {
	u.BalancedID = bal.ID
	u.IsVerified = bal.IsVerified()
}

func (u *User) ConvertToBalanced(data map[string]interface{}) *balanced.Customer {
	bUser := new(balanced.Customer)
	bUser.ID = u.BalancedID
	bUser.Name = data["fullFirstName"].(string) + " " + data["lastName"].(string)
	bUser.Phone = data["phoneNumber"].(string)
	bUser.Email = data["email"].(string)
	bUser.DobYear = data["dobYear"].(int)
	bUser.DobMonth = data["dobMonth"].(int)
	bUser.Address = new(balanced.Address)
	bUser.Address.Line1 = data["streetAddress"].(string)
	bUser.Address.PostalCode = data["postalCode"].(string)
	bUser.Address.CountryCode = "US"
	bUser.SSNLast4 = data["ssnLastFour"].(string)
	return bUser
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
		userBal.ID = user.BalancedID
		balCard := c.ConvertToBalancedCard()
		bErrors := balCard.AssociateWithCustomer(userBal)
		if len(bErrors) > 0 {
			fmt.Println(formatBalancedErrors(bErrors))
		}
		c.ConvertFromBalancedCard(balCard)
		u.Save()
	}(u, card)
	u.Cards = append(u.Cards, card)
	return nil
}

func (u *User) DeleteCard(cardID string) error {
	for i, card := range u.Cards {
		if card.ID == cardID {
			balCard := new(balanced.Card)
			balCard.ID = card.BalancedID
			bErrors := balanced.Delete(balCard)
			if len(bErrors) > 0 {
				return formatBalancedErrors(bErrors)
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
	userBal.ID = u.BalancedID
	balAccount := ba.ConvertToBalancedAccount()
	bErrors := balAccount.AssociateWithCustomer(userBal)
	if len(bErrors) > 0 {
		return formatBalancedErrors(bErrors)
	}
	verification, bErrors := balAccount.Verify()
	if len(bErrors) > 0 {
		return formatBalancedErrors(bErrors)
	}
	ba.VerificationID = verification.ID
	ba.AccountNumber = balAccount.AccountNumber
	ba.RoutingNumber = balAccount.RoutingNumber
	u.Accounts = append(u.Accounts, ba)
	return nil
}

func (u *User) DeleteBankAccount(accountID string) error {
	for i, account := range u.Accounts {
		if account.ID == accountID {
			balAccount := new(balanced.BankAccount)
			balAccount.ID = account.BalancedID
			bErrors := balanced.Delete(balAccount)
			if len(bErrors) > 0 {
				return formatBalancedErrors(bErrors)
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
