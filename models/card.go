package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/balanced-go"
	"strconv"
	// "log"
)

type Card struct {
	ID              string `json:"id"`
	BalancedID      string `json:"balanced_id"`
	LastFour        string `json:"last_four"`
	ExpirationMonth int    `json:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year"`
}
type Cards []*Card

func (c Cards) ToJSON() []byte {
	jsonString := `[`
	for i, card := range c {
		var prefix string
		if i > 0 {
			prefix = `, `
		}
		cardJSON := prefix + `{` +
			`"id":"` + card.ID + `",` +
			`"last_four":"` + card.LastFour + `",` +
			`"expiration_month":` + strconv.Itoa(card.ExpirationMonth) + `,` +
			`"expiration_year":` + strconv.Itoa(card.ExpirationYear) + `}`
		jsonString += cardJSON
	}
	jsonString += `]`
	return []byte(jsonString)
}

func NewCard() *Card {
	id, _ := uuid.NewV4()
	return &Card{
		ID: id.String(),
	}
}

func (c *Card) ConvertFromBalancedCard(balCard *balanced.Card) {
	c.BalancedID = balCard.ID
	cardNumberLength := len(balCard.Number)
	c.LastFour = balCard.Number[cardNumberLength-4 : cardNumberLength]
	c.ExpirationMonth = balCard.ExpirationMonth
	c.ExpirationYear = balCard.ExpirationYear
}

func (c *Card) ConvertToBalancedCard() (balCard *balanced.Card) {
	balCard = new(balanced.Card)
	balCard.ID = c.BalancedID
	return balCard
}
