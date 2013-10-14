package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	// "log"
)

type Card struct {
	ID              string `json:"id"`
	URI             string `json:"uri"`
	LastFour        string `json:"last_four"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
}

func NewCard() *Card {
	id, _ := uuid.NewV4()
	return &Card{
		ID: id.String(),
	}
}

func (c *Card) ConvertBalancedCard(balCard *balanced.Card) {
	c.URI = balCard.URI
	c.LastFour = balCard.LastFour
	c.ExpirationMonth = balCard.ExpirationMonth
	c.ExpirationYear = balCard.ExpirationYear
}

func DeleteCard(userID string, cardID string, ctx *DB.Context) {
	m := make(map[string]interface{})

	change := mgo.Change{
		Update:    bson.M{"$pull": bson.M{"cards": bson.M{"id": cardID}}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("usersbal")
	query := coll.Find(bson.M{
		"_id": userID,
	})

	user := new(User)
	query.One(&user)
	balCard := new(balanced.Card)

	for _, card := range user.Cards {
		if card.ID == cardID {
			balCard.URI = card.URI
			bError := balCard.Delete()
			if bError != nil {
				return
			}

			_, _ = query.Apply(change, &m)
		}
	}

}
