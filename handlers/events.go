package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-PaymentInfo/models"
	"log"
	"net/http"
)

var Connection *amqp.Connection

func dialRMQ() {
	var err error
	Connection, err = amqp.Dial(config.RMQBroker)
	if err != nil {
		panic(err)
	}
}

type Event struct {
	Name string
	Body []byte
}

type Events []*Event

func (e Events) Publish() {
	ch := getChannel()
	defer ch.Close()
	for _, event := range e {
		event.PublishOnChannel(ch)
	}
}

func (e *Event) PublishOnChannel(ch *amqp.Channel) {
	if ch == nil {
		ch = getChannel()
		defer ch.Close()
	}

	ch.ExchangeDeclare(
		config.MainExchange, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	)

	ch.Publish(
		config.MainExchange, // name
		e.Name,              // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        e.Body,
		})
}

func getChannel() *amqp.Channel {
	ch, err := Connection.Channel()
	if ch == nil {
		dialRMQ()
		ch, err = Connection.Channel()
	}
	if err != nil {
		log.Println("rmq", err.Error())
	}
	return ch
}

var BalancedCardType string = "CardBalanced"
var BalancedBankType string = "BankBalanced"

func UpdatePaymentSubmitted(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var message struct {
		PaymentID            string  `json:"paymentID"`
		Amount               float64 `json:"amount"`
		UserID               string  `json:"userID"`
		CreditSourceID       string  `json:"creditSourceID,omitempty"`
		CreditSourceBalanced string  `json:"creditSourceBalanced,omitempty"`
	}

	json.Unmarshal(body, &message)

	user, err := models.FindUserByID(message.UserID)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not find user", err.Error()), http.StatusBadRequest
	}

	bankAccount := user.GetBankAccount(message.CreditSourceID)
	if bankAccount == nil {
		return nil, fmt.Errorf("Error: could not find bank account"), http.StatusBadRequest
	}

	message.CreditSourceBalanced = bankAccount.BalancedID

	j, _ := json.Marshal(message)
	events := Events{&Event{"paymentinfo.credit", j}}
	events.Publish()
	log.Println("update payment submitted", string(j))

	return nil, nil, http.StatusOK
}

func UpdatePaymentAccepted(params map[string]interface{}, body []byte) ([]byte, error, int) {
	fmt.Println("payment accepted")
	var message struct {
		PaymentID           string  `json:"paymentID"`
		Amount              float64 `json:"amount"`
		UserID              string  `json:"userID"`
		DebitSourceID       string  `json:"debitSourceID,omitempty"`
		DebitSourceBalanced string  `json:"debitSourceBalanced,omitempty"`
		PaymentType         string  `json:"paymentType,omitempty"`
	}

	json.Unmarshal(body, &message)

	user, err := models.FindUserByID(message.UserID)
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error: could not find user", err.Error()), http.StatusBadRequest
	}

	bankAccount := user.GetBankAccount(message.DebitSourceID)
	card := user.GetCard(message.DebitSourceID)
	if bankAccount == nil && card == nil {
		return nil, fmt.Errorf("Error: could not find bank account or card"), http.StatusBadRequest
	}

	if bankAccount != nil {
		message.DebitSourceBalanced = bankAccount.BalancedID
		message.PaymentType = "BankBalanced"
	} else if card != nil {
		message.DebitSourceBalanced = card.BalancedID
		message.PaymentType = "CardBalanced"
	}

	j, _ := json.Marshal(message)
	events := Events{&Event{"paymentinfo.debit", j}}
	events.Publish()
	fmt.Println("published")

	return nil, nil, http.StatusOK
}
