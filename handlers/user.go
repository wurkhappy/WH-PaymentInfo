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

	userBal := new(balanced.Customer)
	bError := userBal.Create()

	if bError != nil {
		errorCode, _ := strconv.Atoi(bError.StatusCode)
		http.Error(w, bError.Description, errorCode)
		return
	}

	user := new(models.User)

	user.ID = id
	user.URI = userBal.URI

	err := user.SaveWithCtx(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(`{}`))
}
