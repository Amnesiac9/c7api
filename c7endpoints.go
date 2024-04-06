package c7api

type endpoints struct {
	Auth           string // "https://api.commerce7.com/v1/account/user" - https://developer.commerce7.com/docs/commerce7-apis
	Cart           string // "https://api.commerce7.com/v1/cart" - https://developer.commerce7.com/docs/orders
	Customer       string // "https://api.commerce7.com/v1/customer" - https://developer.commerce7.com/docs/customers
	CustomerId     string // "https://api.commerce7.com/v1/customer/{:id}" - https://developer.commerce7.com/docs/customers
	FulfillmentAll string // "https://api.commerce7.com/v1/order/{:id}/fulfillment/all" - https://developer.commerce7.com/docs/fulfillment
	Fulfillment    string // "https://api.commerce7.com/v1/order/{:id}/fulfillment" - https://developer.commerce7.com/docs/fulfillment
	Inventory      string // "https://api.commerce7.com/v1/inventory" - https://developer.commerce7.com/docs/inventory
	InventoryTrans string // "https://api.commerce7.com/v1/inventory-transaction" - https://developer.commerce7.com/docs/inventory
	Order          string // "https://api.commerce7.com/v1/order" - https://developer.commerce7.com/docs/orders
	Product        string // "https://api.commerce7.com/v1/product" - https://developer.commerce7.com/docs/products
}

func GetEndpointsV1() *endpoints {
	return &endpoints{
		Auth:           "https://api.commerce7.com/v1/account/user",
		Cart:           "https://api.commerce7.com/v1/cart",
		Customer:       "https://api.commerce7.com/v1/customer",
		CustomerId:     "https://api.commerce7.com/v1/customer/{:id}",
		FulfillmentAll: "https://api.commerce7.com/v1/order/{:id}/fulfillment/all",
		Fulfillment:    "https://api.commerce7.com/v1/order/{:id}/fulfillment",
		Inventory:      "https://api.commerce7.com/v1/inventory",
		InventoryTrans: "https://api.commerce7.com/v1/inventory-transaction",
		Order:          "https://api.commerce7.com/v1/order",
		Product:        "https://api.commerce7.com/v1/product",
	}
}
