package main

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/wurkhappy/WH-PaymentInfo/handlers"
)

//order matters so most general should go towards the bottom
var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "/user/:id",
			Dest: map[string]interface{}{
				"POST": handlers.CreateUser,
				"GET":  handlers.GetUser,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/cards",
			Dest: map[string]interface{}{
				"GET":  handlers.GetCards,
				"POST": handlers.SaveCard,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/bank_account",
			Dest: map[string]interface{}{
				"POST": handlers.SaveBankAccount,
				"GET":  handlers.GetBankAccounts,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/cards/:cardID/uri",
			Dest: map[string]interface{}{
				"GET": handlers.GetCardURI,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/cards/:cardID",
			Dest: map[string]interface{}{
				"DELETE": handlers.DeleteCard,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/bank_account/:accountID/uri",
			Dest: map[string]interface{}{
				"GET": handlers.GetBankAccountURI,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/bank_account/:accountID",
			Dest: map[string]interface{}{
				"DELETE": handlers.DeleteBankAccount,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/bank_account/:accountID/verify",
			Dest: map[string]interface{}{
				"POST": handlers.VerifyBankAccount,
			},
		},
	},
}
