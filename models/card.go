package models

import (
	"github.com/wurkhappy/Balanced-go"
)

type Card struct {
	URI             string `json:"uri"`
	LastFour        string `json:"last_four"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
}

func (c *Card) ConvertBalancedCard(balCard *balanced.Card) {
	c.URI = balCard.URI
	c.URI = balCard.LastFour
	c.ExpirationMonth = balCard.ExpirationMonth
	c.ExpirationYear = balCard.ExpirationYear
}
