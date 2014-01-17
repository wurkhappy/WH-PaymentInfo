package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"net/http"
)

func GetUser(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	user, _ := models.FindUserByID(id)

	u, _ := json.Marshal(user)
	return u, nil, http.StatusOK
}

func CreateUser(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)

	user, _ := models.FindUserByID(id)
	if user != nil {
		return []byte(`{}`), nil, http.StatusOK
	}

	user, err := models.CreateUserWithID(id)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not create user", err.Error()), http.StatusBadRequest
	}

	err = user.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not save user", err.Error()), http.StatusBadRequest
	}

	return []byte(`{}`), nil, http.StatusOK
}
