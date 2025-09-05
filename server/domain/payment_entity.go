package domain

// PaymentRequest is the struct for incoming payment API calls.
type PaymentRequest struct {
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	WebhookURL string  `json:"webhookUrl"`
}

// this Payment struct represents the internal structure, this includes the id and status
type Payment struct {
	ID         string  `json:"id"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	WebhookURL string  `json:"-"` // here we exclude the webhookurl to be sent in the respone for security reasons.
}