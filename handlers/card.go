package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"net/http"
	"strconv"
)

func SaveCard(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	card := new(models.Card)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	err := json.Unmarshal(buf.Bytes(), &card)
	if card.ExpirationMonth == 0 {
		var m map[string]interface{}
		json.Unmarshal(buf.Bytes(), &m)
		month, _ := strconv.Atoi(m["expiration_month"].(string))
		year, _ := strconv.Atoi(m["expiration_year"].(string))
		card.ExpirationMonth = month
		card.ExpirationYear = year
	}

	cardID, _ := uuid.NewV4()
	card.ID = cardID.String()

	user, err := models.FindUserByID(id)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	err = user.AddCreditCard(card)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = user.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, _ := json.Marshal(card)
	w.Write(u)
}

func GetCards(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]

	user, err := models.FindUserByID(id)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	jsonBytes, _ := json.Marshal(user.Cards)
	w.Write(jsonBytes)
}

func DeleteCard(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	cardID := vars["cardID"]

	user, err := models.FindUserByID(id)
	if err != nil {
		http.Error(w, "Error: couldn't find user", http.StatusBadRequest)
		return
	}

	err = user.DeleteCard(cardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
