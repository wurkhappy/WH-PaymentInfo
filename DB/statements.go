package DB

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	// "log"
)

var UpsertUser *sql.Stmt
var FindUserByEmail *sql.Stmt
var FindUserByID *sql.Stmt
var DeleteUser *sql.Stmt
var FindUsers *sql.Stmt
var SyncWithExistingInvitation *sql.Stmt

func CreateStatements() {
	var err error

	UpsertUser, err = DB.Prepare("SELECT upsert_balanceduser($1, $2)")
	if err != nil {
		panic(err)
	}

	FindUserByID, err = DB.Prepare("SELECT data FROM balanced_user WHERE id = $1")
	if err != nil {
		panic(err)
	}

	FindUserByEmail, err = DB.Prepare("SELECT password, data FROM wh_user WHERE data->>'email' = $1")
	if err != nil {
		panic(err)
	}

	DeleteUser, err = DB.Prepare("DELETE FROM wh_user WHERE id = $1")
	if err != nil {
		panic(err)
	}

	FindUsers, err = DB.Prepare("SELECT data FROM wh_user WHERE id = ANY($1)")
	if err != nil {
		panic(err)
	}
}
