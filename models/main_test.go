package models

import (
	"flag"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-PaymentInfo/DB"
	// "log"
)

var balMarketplaceID string = "TEST-MP1f775iSL82BucxjmR83cOk"
var useBal = flag.Bool("balanced", false, "run tests with balanced integration")

func init() {
	flag.Parse()
	balanced.Username = config.BalancedUsername
	DB.Name = "testdb"
	DB.Setup()
	DB.CreateStatements()
}
