package models

type BankAccount struct {
	AccountNumber    string `json:"account_number,omitempty"`
	URI              string `json:"uri,omitempty"`
	RoutingNumber    string `json:"routing_number,omitempty"`
	VerificationsURI string `json:"verifications_uri,omitempty"`
	VerificationURI  string `json:"verification_uri,omitempty"`
	CreditsURI       string `json:"credits_uri,omitempty"`
}
