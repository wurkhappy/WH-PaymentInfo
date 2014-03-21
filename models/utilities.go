package models

import (
	"fmt"
	"github.com/wurkhappy/balanced-go"
)

func formatBalancedErrors(bErrors []*balanced.BalancedError) error {
	var errorMessages string
	for _, bError := range bErrors {
		errorMessages += bError.Description + ", "
	}
	return fmt.Errorf(errorMessages)
}