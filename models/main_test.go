package models

import (
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	"flag"
	// "log"
)

var balMarketplaceID string = "TEST-MP1f775iSL82BucxjmR83cOk"
var useBal = flag.Bool("balanced", false, "run tests with balanced integration")
func init() {
	flag.Parse()
	balanced.Username = "ak-test-x9PqPQUtpvUtnXsZqBL4rXGAE8WvvqoJ"
	DB.Name = "testdb"
	DB.Setup()
	DB.CreateStatements()
}
