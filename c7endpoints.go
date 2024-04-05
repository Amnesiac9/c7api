package c7api

type Endpoint struct {
	Auth     string
	Customer string
	Order    string
}

func GetEndpointsV1() *Endpoint {
	endpoint := Endpoint{
		Auth:     "https://api.commerce7.com/v1/account/user",
		Customer: "",
		Order:    "",
	}

	return &endpoint
}
