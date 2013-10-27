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
