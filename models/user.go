package models

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"log"
)

type User struct {
	ID         string       `json:"id" bson:"_id"`
	URI        string       `json:"uri"`
	DebitsURI  string       `json:"debitsURI"`
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
	bError := userBal.Create()
	if bError != nil {
		return nil, fmt.Errorf("%s", bError.Description)
	}
	u.ConvertFromBalanced(userBal)
	return u, nil
}

func (u *User) UpdateWithMap(m map[string]interface{}) error {
	bUser := u.ConvertToBalanced(m)
	bError := bUser.Update()
	if bError != nil {
		return fmt.Errorf("%s", bError.Description)
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
	u.URI = bal.URI
	u.DebitsURI = bal.DebitsURI
	u.IsVerified = bal.IdentityVerified
}

func (u *User) ConvertToBalanced(data map[string]interface{}) *balanced.Customer {
	bUser := new(balanced.Customer)
	bUser.URI = u.URI
	bUser.Name = data["fullFirstName"].(string) + " " + data["lastName"].(string)
	bUser.Phone = data["phoneNumber"].(string)
	bUser.Email = data["email"].(string)
	year := data["dobYear"].(string)
	month := data["dobMonth"].(string)
	var monthString string
	if len(month) == 1 {
		monthString = "0"
	}
	monthString += month
	bUser.Dob = year + "-" + monthString
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
