package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/Balanced-go"
	// "log"
)

type Card struct {
	ID              string `json:"id"`
	URI             string `json:"uri"`
	LastFour        string `json:"last_four"`
	ExpirationMonth int    `json:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year"`
}

func NewCard() *Card {
	id, _ := uuid.NewV4()
	return &Card{
		ID: id.String(),
	}
}

func (c *Card) ConvertFromBalancedCard(balCard *balanced.Card) {
	c.URI = balCard.URI
	c.LastFour = balCard.LastFour
	c.ExpirationMonth = balCard.ExpirationMonth
	c.ExpirationYear = balCard.ExpirationYear
}

func (c *Card) ConvertToBalancedCard() (balCard *balanced.Card) {
	balCard = new(balanced.Card)
	balCard.URI = c.URI
	balCard.LastFour = c.LastFour
	balCard.ExpirationMonth = c.ExpirationMonth
	balCard.ExpirationYear = c.ExpirationYear
	return balCard
}
