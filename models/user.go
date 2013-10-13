package models

import (
	// "github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"labix.org/v2/mgo/bson"
	"log"
)

type User struct {
	ID       string     `json:"id" bson:"_id"`
	URI      string     `json:"id"`
	Cards    []*Card    `json:"cards"`
	Accounts []*Account `json:"accounts"`
}

func FindUserByID(id string, ctx *DB.Context) (u *User, err error) {
	err = ctx.Database.C("usersbal").Find(bson.M{"_id": id}).One(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) SaveWithCtx(ctx *DB.Context) (err error) {
	coll := ctx.Database.C("usersbal")
	if _, err := coll.UpsertId(u.ID, &u); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
