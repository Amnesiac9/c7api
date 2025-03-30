package c7api

type endpoints struct {
	Auth                 string // "https://api.commerce7.com/v1/account/user" - https://developer.commerce7.com/docs/commerce7-apis
	Cart                 string // "https://api.commerce7.com/v1/cart" - https://developer.commerce7.com/docs/orders
	Customer             string // "https://api.commerce7.com/v1/customer" - https://developer.commerce7.com/docs/customers
	CustomerId           string // "https://api.commerce7.com/v1/customer/{:id}" - https://developer.commerce7.com/docs/customers
	FulfillmentAll       string // "https://api.commerce7.com/v1/order/{:id}/fulfillment/all" - https://developer.commerce7.com/docs/fulfillment
	Fulfillment          string // "https://api.commerce7.com/v1/order/{:id}/fulfillment" - https://developer.commerce7.com/docs/fulfillment
	Inventory            string // "https://api.commerce7.com/v1/inventory" - https://developer.commerce7.com/docs/inventory
	InventoryTransaction string // "https://api.commerce7.com/v1/inventory-transaction" - https://developer.commerce7.com/docs/inventory
	Order                string // "https://api.commerce7.com/v1/order" - https://developer.commerce7.com/docs/orders
	Product              string // "https://api.commerce7.com/v1/product" - https://developer.commerce7.com/docs/products
	Tag                  string // "https://api.commerce7.com/v1/tag" - example: https://api.commerce7.com/v1/tag/customer?q=&type=Manual - Use q to search for specific keywords, type manual to get only manually applied tags
	LoyaltyTransaction   string // "https://api.commerce7.com/v1/loyalty-transaction" - requires Commerce7 Loyalty Extension https://documentation.commerce7.com/loyalty-feature-overview
	Vendor               string // "https://api.commerce7.com/v1/vendor" - example: https://api.commerce7.com/v1/vendor?q=myvendor
	Collection           string // "https://api.commerce7.com/v1/collection" - example: https://api.commerce7.com/v1/collection?type=Manual
	Department           string // "https://api.commerce7.com/v1/department" - example: https://api.commerce7.com/v1/department?q=wine
}

func GetEndpointsV1() *endpoints {
	return &endpoints{
		Auth:                 "https://api.commerce7.com/v1/account/user",
		Cart:                 "https://api.commerce7.com/v1/cart",
		Customer:             "https://api.commerce7.com/v1/customer",
		CustomerId:           "https://api.commerce7.com/v1/customer/{:id}",
		FulfillmentAll:       "https://api.commerce7.com/v1/order/{:id}/fulfillment/all",
		Fulfillment:          "https://api.commerce7.com/v1/order/{:id}/fulfillment",
		Inventory:            "https://api.commerce7.com/v1/inventory",
		InventoryTransaction: "https://api.commerce7.com/v1/inventory-transaction",
		Order:                "https://api.commerce7.com/v1/order",
		Product:              "https://api.commerce7.com/v1/product",
		Tag:                  "https://api.commerce7.com/v1/tag",
		LoyaltyTransaction:   "https://api.commerce7.com/v1/loyalty-transaction",
		Vendor:               "https://api.commerce7.com/v1/vendor",
		Collection:           "https://api.commerce7.com/v1/collection",
		Department:           "https://api.commerce7.com/v1/department",
	}
}

var Endpoints = GetEndpointsV1()
