package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"testing"
)

func Test_AccountUnit(t *testing.T) {

}

func Test_AccountIntegration(t *testing.T) {
	if !testing.Short() {

		if *useBal {
			test_ConfirmVerification(t)
		}

		DB.DB.Exec("DELETE from balancedUser")
	}
}
func test_ConfirmVerification(t *testing.T) {
	id, _ := uuid.NewV4()
	user, _ := CreateUserWithID(id.String())
	balAccount := createBalancedBankAccount()
	account := new(BankAccount)
	accountid, _ := uuid.NewV4()
	account.ID = accountid.String()
	account.ConvertFromBalancedAccount(balAccount)

	user.AddBankAccount(account)

	err := account.ConfirmVerification(1, 1)
	if err != nil {
		t.Errorf("test_ConfirmVerification--- error confirming bank account %s", err.Error())
	}

	if !account.CanDebit {
		t.Errorf("%s--- can debit not set to true", "test_ConfirmVerification")
	}
}
