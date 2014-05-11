package main

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/wurkhappy/WH-PaymentInfo/handlers"
)

var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "payment.submitted",
			Dest:    handlers.UpdatePaymentSubmitted,
		},
		urlrouter.Route{
			PathExp: "payment.accepted",
			Dest:    handlers.UpdatePaymentAccepted,
		},
		urlrouter.Route{
			PathExp: "user.created",
			Dest:    handlers.CreateUser,
		},
		urlrouter.Route{
			PathExp: "agreement.submitted",
			Dest:    handlers.AgreementSubmitted,
		},
	},
}
