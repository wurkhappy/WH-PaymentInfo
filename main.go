package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"github.com/wurkhappy/WH-PaymentInfo/handlers"
	"labix.org/v2/mgo"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", req.URL.Path[1:])
}

func main() {
	var err error
	DB.Session, err = mgo.Dial(DB.Config["DBURL"])
	if err != nil {
		panic(err)
	}
	balanced.Username = "ak-test-3TJ2ufbpEUv1BDgxpoCstfj5h2XWMAqG"
	r := mux.NewRouter()
	r.Handle("/user/{id}", dbContextMixIn(handlers.GetUser)).Methods("GET")
	r.Handle("/user/{id}", dbContextMixIn(handlers.CreateUser)).Methods("POST")
	r.Handle("/user/{id}/cards", dbContextMixIn(handlers.SaveCard)).Methods("POST")
	http.Handle("/", r)

	http.ListenAndServe(":3120", nil)
}

type dbContextMixIn func(http.ResponseWriter, *http.Request, *DB.Context)

func (h dbContextMixIn) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//create the context
	ctx, err := DB.NewContext(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.Close()

	h(w, req, ctx)
}
