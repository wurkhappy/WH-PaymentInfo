package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"net/http"
	// "log"
)

func SaveBankAccount(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)

	bankAccount := new(models.BankAccount)
	json.Unmarshal(body, &bankAccount)

	accountID, _ := uuid.NewV4()
	bankAccount.ID = accountID.String()

	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error: could not find user"), http.StatusBadRequest
	}

	err = user.AddBankAccount(bankAccount)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Could not add bank account: ", err.Error()), http.StatusBadRequest
	}

	err = user.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Could not save user: ", err.Error()), http.StatusBadRequest
	}

	ba, _ := json.Marshal(bankAccount)
	return ba, nil, http.StatusOK
}

func GetBankAccounts(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)

	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error: could not find user"), http.StatusBadRequest
	}

	jsonBytes := user.Accounts.ToJSON()
	return jsonBytes, nil, http.StatusOK
}

func DeleteBankAccount(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	accountID := params["accountID"].(string)

	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error: could not find user"), http.StatusBadRequest
	}

	err = user.DeleteBankAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not delete account", err.Error()), http.StatusBadRequest
	}
	return nil, nil, http.StatusOK
}

func VerifyBankAccount(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	accountID := params["accountID"].(string)

	var amounts struct {
		Amount1 float64 `json:"amount_1"`
		Amount2 float64 `json:"amount_2"`
	}

	json.Unmarshal(body, &amounts)

	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not find user", err.Error()), http.StatusBadRequest
	}

	bankAccount := user.GetBankAccount(accountID)
	err = bankAccount.ConfirmVerification(amounts.Amount1, amounts.Amount2)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not confirm verification", err.Error()), http.StatusBadRequest
	}

	err = user.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not save user", err.Error()), http.StatusBadRequest
	}
	return nil, nil, http.StatusOK
}

func GetBankAccountURI(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	accountID := params["accountID"].(string)

	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not find user", err.Error()), http.StatusBadRequest
	}

	bankAccount := user.GetBankAccount(accountID)
	if bankAccount == nil {
		return nil, fmt.Errorf("Error: could not find bank account"), http.StatusBadRequest
	}
	jsonBytes, _ := json.Marshal(bankAccount)
	return jsonBytes, nil, http.StatusOK
}
