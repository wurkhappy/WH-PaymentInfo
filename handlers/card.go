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

func SaveCard(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	vars := mux.Vars(req)
	id := vars["id"]

	balCard := new(balanced.Card)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &balCard)

	user, err := models.FindUserByID(id, ctx)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	userBal := new(balanced.Customer)
	userBal.URI = user.URI
	bError := userBal.AddCreditCard(balCard)
	if bError != nil {
		errorCode, _ := strconv.Atoi(bError.StatusCode)
		http.Error(w, bError.Description, errorCode)
		return
	}

	card := new(models.Card)
	card.ConvertBalancedCard(balCard)

	user.Cards = append(user.Cards, card)

	err = user.SaveWithCtx(ctx)
		if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, _ := json.Marshal(card)
	w.Write(u)
}
