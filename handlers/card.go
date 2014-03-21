package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"net/http"
)

func SaveCard(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)

	card := new(models.Card)
	json.Unmarshal(body, &card)

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

	jsonBytes := user.Cards.ToJSON()
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

func GetCardURI(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	cardID := params["cardID"].(string)

	user, err := models.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not find user", err.Error()), http.StatusBadRequest
	}

	card := user.GetCard(cardID)
	if card == nil {
		return nil, fmt.Errorf("Error: could not find card"), http.StatusBadRequest
	}
	jsonBytes, _ := json.Marshal(card)
	return jsonBytes, nil, http.StatusOK
}
