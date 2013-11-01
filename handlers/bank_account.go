package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"net/http"
	"strconv"
	// "log"
)

func SaveBankAccount(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	vars := mux.Vars(req)
	id := vars["id"]

	balAccount := new(balanced.BankAccount)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &balAccount)

	user, err := models.FindUserByID(id, ctx)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	userBal := new(balanced.Customer)
	userBal.URI = user.URI
	bError := userBal.AddBankAccount(balAccount)
	if bError != nil {
		errorCode, _ := strconv.Atoi(bError.StatusCode)
		http.Error(w, bError.Description, errorCode)
		return
	}
	_, bError = balAccount.Verify()
	if bError != nil {
		errorCode, _ := strconv.Atoi(bError.StatusCode)
		http.Error(w, bError.Description, errorCode)
		return
	}

	bankAccount := models.NewBankAccount()
	bankAccount.ConvertBalancedAccount(balAccount)

	user.Accounts = append(user.Accounts, bankAccount)

	err = user.SaveWithCtx(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, _ := json.Marshal(bankAccount)
	w.Write(u)
}

func GetBankAccounts(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	vars := mux.Vars(req)
	id := vars["id"]

	user, err := models.FindUserByID(id, ctx)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	jsonBytes, _ := json.Marshal(user.Accounts)
	w.Write(jsonBytes)
}

func DeleteBankAccount(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	vars := mux.Vars(req)
	id := vars["id"]
	accountID := vars["accountID"]

	models.DeleteBankAccount(id, accountID, ctx)
}

func VerifyBankAccount(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	vars := mux.Vars(req)
	id := vars["id"]
	accountID := vars["accountID"]

	var amounts struct {
		Amount1 float64 `json:"amount_1"`
		Amount2 float64 `json:"amount_2"`
	}

	user, err := models.FindUserByID(id, ctx)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	var bankAccount *models.BankAccount
	for _, account := range user.Accounts {
		if account.ID == accountID {
			bankAccount = account
		}
	}
	balAccount := new(balanced.BankAccount)
	balAccount.VerificationURI = bankAccount.VerificationURI

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &amounts)

	verification, bError := balAccount.ConfirmVerification(amounts.Amount1, amounts.Amount2)
	if bError != nil {
		errorCode, _ := strconv.Atoi(bError.StatusCode)
		http.Error(w, bError.Description, errorCode)
		return
	}
	if verification.State == "verified" {
		bankAccount.CanDebit = true
		err = user.SaveWithCtx(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Account not verified", http.StatusBadRequest)
		return
	}
}
