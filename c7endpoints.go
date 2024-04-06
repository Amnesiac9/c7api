package c7api

// var Endpoints struct {
// 	Base         string // "https://api.commerce7.com/v1"
// 	Auth         string // "https://api.commerce7.com/v1/account/user"
// 	Cart         string // "https://api.commerce7.com/v1/cart"
// 	Customer     string // "https://api.commerce7.com/v1/customer"
// 	CustomerId   string // "https://api.commerce7.com/v1/customer/{:id}"
// 	CustomerAddr string // "https://api.commerce7.com/v1/customer-address"
// 	Order        string
// } {}

const (
	Endpoint               string = "https://api.commerce7.com/v1"                             // https://developer.commerce7.com/docs/commerce7-apis
	EndpointAuth           string = "https://api.commerce7.com/v1/account/user"                // https://developer.commerce7.com/docs/commerce7-apis
	EndpointCart           string = "https://api.commerce7.com/v1/cart"                        // https://developer.commerce7.com/docs/orders
	EndpointCustomer       string = "https://api.commerce7.com/v1/customer"                    // https://developer.commerce7.com/docs/customers
	EndpointCustomerId     string = "https://api.commerce7.com/v1/customer/{:id}"              // https://developer.commerce7.com/docs/customers
	EndpointFulfillmentAll string = "https://api.commerce7.com/v1/order/{:id}/fulfillment/all" // https://developer.commerce7.com/docs/fulfillment
	EndpointFulfillment    string = "https://api.commerce7.com/v1/order/{:id}/fulfillment"     // https://developer.commerce7.com/docs/fulfillment
	EndpointInventory      string = "https://api.commerce7.com/v1/inventory"                   // https://developer.commerce7.com/docs/inventory
	EndpointInventoryTrans string = "https://api.commerce7.com/v1/inventory-transaction"       // https://developer.commerce7.com/docs/inventory
	EndpointOrder          string = "https://api.commerce7.com/v1/order"                       // https://developer.commerce7.com/docs/orders
	EndpointProduct        string = "https://api.commerce7.com/v1/product"                     // https://developer.commerce7.com/docs/products
)
