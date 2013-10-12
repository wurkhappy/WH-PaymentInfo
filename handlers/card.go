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

func SaveCardUri(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	vars := mux.Vars(req)
	id := vars["id"]

	card := new(balanced.Card)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &card)

	user, err := models.FindUserByID(id, ctx)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}

	userBal := new(balanced.Customer)
	userBal.URI = user.URI
	bError := userBal.AddCreditCard(card)
	if bError != nil {
		errorCode, _ := strconv.Atoi(bError.StatusCode)
		http.Error(w, "Error", errorCode)
	}

	user.Cards = append(user.Cards)

	user.SaveWithCtx(ctx)

	u, _ := json.Marshal(user)
	w.Write(u)
}
