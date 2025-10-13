package c7api

type GiftCardPost struct {
	Title         string `json:"title"`
	Type          string `json:"type"`
	Code          string `json:"code,omitempty"` // optional
	Status        string `json:"status"`
	InitialAmount int    `json:"initialAmount"`
	Notes         string `json:"notes"`
	ExpiryDate    string `json:"expiryDate"` // or time.Time if you plan to parse it
}
