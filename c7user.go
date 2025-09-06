package c7api

type NewUser struct {
	Tenant string `json:"tenantId"`
	User   struct {
		C7UserID  string `json:"id"`
		FirstName string `json:"firstName" form:"firstName"`
		LastName  string `json:"lastName" form:"lastName"`
		Email     string `json:"email" form:"email"`
	} `json:"user"`
}
