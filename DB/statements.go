package DB

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	// "log"
)

var UpsertUser *sql.Stmt
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
}
