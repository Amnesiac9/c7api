package c7api

type NewUser struct {
	TenantID string `json:"tenantId"`
	User     struct {
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	} `json:"user"`
}
