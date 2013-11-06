package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"net/http"
	"strconv"
	// "bytes"
)

func GetUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	vars := mux.Vars(req)
	id := vars["id"]
	user, _ := models.FindUserByID(id, ctx)

	u, _ := json.Marshal(user)
	w.Write(u)
}

func CreateUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
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
