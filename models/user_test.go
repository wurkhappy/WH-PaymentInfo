package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"github.com/wurkhappy/balanced-go"
	"log"
	"testing"
	"time"
)

func Test_UserUnit(t *testing.T) {
	test_ConvertFromBalanced(t)
}

func Test_UserIntegration(t *testing.T) {
	if !testing.Short() {
		testSaveUser(t)
		testFindUserByID(t)

		if *useBal {
			testCreateUserWithID(t)
			test_AddCreditCard(t)
			test_DeleteCard(t)
			test_AddBankAccount(t)
			test_DeleteBankAccount(t)
		}

		DB.DB.Exec("DELETE from balancedUser")
	}
}

func testSaveUser(t *testing.T) {
	user := new(User)
	id, _ := uuid.NewV4()
	user.ID = id.String()
	err := user.Save()

	if err != nil {
		t.Errorf("%s--- error is:%v", "testSaveUser", err)
	}
}

func testFindUserByID(t *testing.T) {
	user := new(User)
	id, _ := uuid.NewV4()
	user.ID = id.String()
	user.Save()

	u, err := FindUserByID(user.ID)
	if err != nil {
		t.Errorf("testFindUserByID--- error finding user %v", err)
	}

	if u == nil {
		t.Errorf("%s--- user was not found", "testFindUserByID")
	}

	_, err = FindUserByID("invalidID")
	if err == nil {
		t.Errorf("%s--- DB returned a bad request", "testFindUserByID")
	}
}

func testCreateUserWithID(t *testing.T) {
	id, _ := uuid.NewV4()
	user, err := CreateUserWithID(id.String())
	if err != nil {
		t.Errorf("testCreateUserWithID--- error creating user %v", err)
	}

	if user.URI == "" {
		t.Errorf("%s--- user doesn't have a uri", "testCreateUserWithID")
	}

	if user.DebitsURI == "" {
		t.Errorf("%s--- user doesn't have a debits uri", "testCreateUserWithID")
	}

}

func test_ConvertFromBalanced(t *testing.T) {
	bal := new(balanced.Customer)
	bal.URI = "test.com"
	bal.DebitsURI = "test/debits"

	user := new(User)
	user.ConvertFromBalanced(bal)
}

func test_AddCreditCard(t *testing.T) {
	id, _ := uuid.NewV4()
	user, _ := CreateUserWithID(id.String())
	balCard := createBalancedCard()
	card := new(Card)
	cardid, _ := uuid.NewV4()
	card.ID = cardid.String()
	card.ConvertFromBalancedCard(balCard)

	err := user.AddCreditCard(card)
	if err != nil {
		t.Errorf("test_AddCreditCard--- error adding card %s", err.Error())
	}

	if len(user.Cards) != 1 {
		t.Errorf("%s--- wrong number of cards returned", "test_AddCreditCard")
	}

}

func test_DeleteCard(t *testing.T) {
	id, _ := uuid.NewV4()
	user, _ := CreateUserWithID(id.String())
	balCard := createBalancedCard()
	card := new(Card)
	cardid, _ := uuid.NewV4()
	card.ID = cardid.String()
	card.ConvertFromBalancedCard(balCard)

	err := user.AddCreditCard(card)
	if err != nil {
		t.Errorf("test_DeleteCard--- error adding card %s", err.Error())
	}

	err = user.DeleteCard(card.ID)
	if err != nil {
		t.Errorf("test_DeleteCard--- error deleting card %s", err.Error())
	}
	if len(user.Cards) != 0 {
		t.Errorf("%s--- wrong number of cards returned", "test_AddCreditCard")
	}

}

func test_AddBankAccount(t *testing.T) {
	id, _ := uuid.NewV4()
	user, _ := CreateUserWithID(id.String())
	balAccount := createBalancedBankAccount()
	account := new(BankAccount)
	accountid, _ := uuid.NewV4()
	account.ID = accountid.String()
	account.ConvertFromBalancedAccount(balAccount)

	err := user.AddBankAccount(account)
	if err != nil {
		t.Errorf("test_AddBankAccount--- error adding bank account %s", err.Error())
	}

	if len(user.Accounts) != 1 {
		t.Errorf("%s--- wrong number of bank accounts returned", "test_AddBankAccount")
	}

}

func test_DeleteBankAccount(t *testing.T) {
	id, _ := uuid.NewV4()
	user, _ := CreateUserWithID(id.String())
	balAccount := createBalancedBankAccount()
	account := new(BankAccount)
	accountid, _ := uuid.NewV4()
	account.ID = accountid.String()
	account.ConvertFromBalancedAccount(balAccount)

	err := user.AddBankAccount(account)
	if err != nil {
		t.Errorf("test_DeleteBankAccount--- error adding bank account %s", err.Error())
	}

	err = user.DeleteBankAccount(account.ID)
	if err != nil {
		t.Errorf("test_DeleteBankAccount--- error deleting bank account %s", err.Error())
	}
	if len(user.Accounts) != 0 {
		t.Errorf("%s--- wrong number of bank accounts returned", "test_DeleteBankAccount")
	}

}

func createBalancedCard() *balanced.Card {
	balCard := new(balanced.Card)
	balCard.CardNumber = "4111111111111111"
	balCard.ExpirationYear = time.Now().Year() + 1
	balCard.ExpirationMonth = 1
	balCard.SecurityCode = "123"

	bError := balCard.Create(balMarketplaceID)
	if bError != nil {
		log.Printf("createBalancedCard err:  %v", bError)
	}

	return balCard
}

func createBalancedBankAccount() *balanced.BankAccount {
	balAccount := new(balanced.BankAccount)
	balAccount.RoutingNumber = "021000021"
	balAccount.Type = "checking"
	balAccount.AccountNumber = "9900000002"
	balAccount.Name = "Wurk Happy"

	bError := balAccount.Create()
	if bError != nil {
		log.Printf("createBalancedBankAccount err:  %v", bError)
	}

	return balAccount
}
