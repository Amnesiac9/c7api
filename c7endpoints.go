package c7api

const API_BASE_URL = "https://api.commerce7.com/v1"

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
	TagXObject           string // https://api.commerce7.com/v1/tag-x-object/{:obj} - Requires an object appended, either customer or order
	LoyaltyTransaction   string // "https://api.commerce7.com/v1/loyalty-transaction" - requires Commerce7 Loyalty Extension https://documentation.commerce7.com/loyalty-feature-overview
	Vendor               string // "https://api.commerce7.com/v1/vendor" - example: https://api.commerce7.com/v1/vendor?q=myvendor
	Collections          string // "https://api.commerce7.com/v1/collection" - example: https://api.commerce7.com/v1/collection?type=Manual
	Department           string // "https://api.commerce7.com/v1/department" - example: https://api.commerce7.com/v1/department?q=wine
	Setting              string // "https://api.commerce7.com/v1/setting" - Used for grabbing the winery info/settings from a given tenant.
}

func GetEndpointsV1() *endpoints {
	return &endpoints{
		Auth:                 API_BASE_URL + "/account/user",
		Cart:                 API_BASE_URL + "/cart",
		Customer:             API_BASE_URL + "/customer",
		CustomerId:           API_BASE_URL + "/customer/{:id}",
		FulfillmentAll:       API_BASE_URL + "/order/{:id}/fulfillment/all",
		Fulfillment:          API_BASE_URL + "/order/{:id}/fulfillment",
		Inventory:            API_BASE_URL + "/inventory",
		InventoryTransaction: API_BASE_URL + "/inventory-transaction",
		Order:                API_BASE_URL + "/order",
		Product:              API_BASE_URL + "/product",
		Tag:                  API_BASE_URL + "/tag",
		TagXObject:           API_BASE_URL + "/tag-x-object/{:obj}",
		LoyaltyTransaction:   API_BASE_URL + "/loyalty-transaction",
		Vendor:               API_BASE_URL + "/vendor",
		Collections:          API_BASE_URL + "/collection",
		Department:           API_BASE_URL + "/department",
		Setting:              API_BASE_URL + "/setting",
	}
}

var Endpoints = GetEndpointsV1()
