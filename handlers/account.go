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

func SaveAccountUri(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	vars := mux.Vars(req)
	id := vars["id"]

	balAccount := new(balanced.BankAccount)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &balAccount)

	user, err := models.FindUserByID(id, ctx)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}

	userBal := new(balanced.Customer)
	userBal.URI = user.URI
	bError := userBal.AddBankAccount(balAccount)
	if bError != nil {
		errorCode, _ := strconv.Atoi(bError.StatusCode)
		http.Error(w, "Error", errorCode)
	}

	account := new(models.BankAccount)
	account.ConvertBalancedAccount(balAccount)

	user.Accounts = append(user.Accounts, account)

	user.SaveWithCtx(ctx)

	u, _ := json.Marshal(user)
	w.Write(u)
}
