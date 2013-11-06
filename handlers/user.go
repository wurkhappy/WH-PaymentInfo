package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"net/http"
	// "bytes"
)

func GetUser(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]
	user, _ := models.FindUserByID(id)

	u, _ := json.Marshal(user)
	w.Write(u)
}

func CreateUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	user, err := models.CreateUserWithID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = user.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(`{}`))
}
