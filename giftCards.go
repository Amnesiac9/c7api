package c7api

// For posting a new gift card
type GiftCardPost struct {
	Title         string `json:"title"`
	Type          string `json:"type"`
	Code          string `json:"code,omitempty"` // optional
	Status        string `json:"status"`
	InitialAmount int    `json:"initialAmount"`
	Notes         string `json:"notes,omitempty"`
	ExpiryDate    string `json:"expiryDate,omitempty"` // optional
}

// Transaction struct for incrementing or decrementing a gift card
type GiftCardTransactionPost struct {
	Amount     int    `json:"amount"`     // Amount in cents
	GiftCardId string `json:"giftCardId"` // UUID of the gift card
}
