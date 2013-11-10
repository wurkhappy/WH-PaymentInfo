package handlers

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"net/http"
	"strconv"
	"fmt"
)

func SaveCard(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)

	card := new(models.Card)
	json.Unmarshal(body, &card)
	if card.ExpirationMonth == 0 {
		var m map[string]interface{}
		json.Unmarshal(body, &m)
		month, _ := strconv.Atoi(m["expiration_month"].(string))
		year, _ := strconv.Atoi(m["expiration_year"].(string))
		card.ExpirationMonth = month
		card.ExpirationYear = year
	}

	cardID, _ := uuid.NewV4()
	card.ID = cardID.String()

	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not find user", err.Error()), http.StatusBadRequest
	}

	err = user.AddCreditCard(card)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not add credit card", err.Error()), http.StatusBadRequest
	}

	err = user.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not save user", err.Error()), http.StatusBadRequest
	}

	u, _ := json.Marshal(card)
	return u, nil, http.StatusOK
}

func GetCards(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)

	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not find user", err.Error()), http.StatusBadRequest
	}

	jsonBytes, _ := json.Marshal(user.Cards)
	return jsonBytes, nil, http.StatusOK
}

func DeleteCard(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	cardID := params["cardID"].(string)

	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not find user", err.Error()), http.StatusBadRequest
	}

	err = user.DeleteCard(cardID)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not delete card", err.Error()), http.StatusBadRequest
	}
	return nil, nil, http.StatusOK
}
