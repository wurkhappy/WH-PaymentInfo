package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"net/http"
	// "log"
)

func SaveBankAccount(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]

	bankAccount := new(models.BankAccount)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &bankAccount)

	accountID, _ := uuid.NewV4()
	bankAccount.ID = accountID.String()

	user, err := models.FindUserByID(id)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	err = user.AddBankAccount(bankAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = user.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, _ := json.Marshal(bankAccount)
	w.Write(u)
}

func GetBankAccounts(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]

	user, err := models.FindUserByID(id)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	jsonBytes, _ := json.Marshal(user.Accounts)
	w.Write(jsonBytes)
}

func DeleteBankAccount(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	accountID := vars["accountID"]

	user, err := models.FindUserByID(id)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	err = user.DeleteBankAccount(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func VerifyBankAccount(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	accountID := vars["accountID"]

	var amounts struct {
		Amount1 float64 `json:"amount_1"`
		Amount2 float64 `json:"amount_2"`
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &amounts)

	user, err := models.FindUserByID(id)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	bankAccount := user.GetBankAccount(accountID)
	err = bankAccount.ConfirmVerification(amounts.Amount1, amounts.Amount2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = user.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
