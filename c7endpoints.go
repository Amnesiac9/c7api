package c7api

const API_URL = "https://api.commerce7.com/v1"
const API_URL_V2 = "https://api.commerce7.com/v2"

type endpoints struct {
	Auth                 string // "https://api.commerce7.com/v1/account/user" - https://developer.commerce7.com/docs/commerce7-apis
	Cart                 string // "https://api.commerce7.com/v1/cart" - https://developer.commerce7.com/docs/orders
	Customer             string // "https://api.commerce7.com/v1/customer" - https://developer.commerce7.com/docs/customers
	CustomerId           string // "https://api.commerce7.com/v1/customer/{:id}" - https://developer.commerce7.com/docs/customers
	ClubMembership       string // "https://api.commerce7.com/v1/club-membership" -
	FulfillmentAll       string // "https://api.commerce7.com/v1/order/{:id}/fulfillment/all" - https://developer.commerce7.com/docs/fulfillment
	Fulfillment          string // "https://api.commerce7.com/v1/order/{:id}/fulfillment" - https://developer.commerce7.com/docs/fulfillment
	GiftCard             string // "https://api.commerce7.com/v1/gift-card" - https://developer.commerce7.com/docs/gift-cards
	GiftCardTransaction  string // "https://api.commerce7.com/v1/gift-card-transaction"
	Inventory            string // "https://api.commerce7.com/v1/inventory" - https://developer.commerce7.com/docs/inventory
	InventoryTransaction string // "https://api.commerce7.com/v1/inventory-transaction" - https://developer.commerce7.com/docs/inventory
	MetaDataConfig       string // "https://api.commerce7.com/v1/meta-data-config/" + "{:obj}" - Requires an object appended; allocation, club-membership, collection, customer, customer-address, order, product, reservation, experience
	MetaDataConfigObj    string // "https://api.commerce7.com/v1/meta-data-config/{:obj}" - Requires an object appended; allocation, club-membership, collection, customer, customer-address, order, product, reservation, experience
	Order                string // "https://api.commerce7.com/v1/order" - https://developer.commerce7.com/docs/orders
	Product              string // "https://api.commerce7.com/v1/product" - https://developer.commerce7.com/docs/products
	Tag                  string // "https://api.commerce7.com/v1/tag" - example: https://api.commerce7.com/v1/tag/customer?q=&type=Manual - Use q to search for specific keywords, type manual to get only manually applied tags
	TagXObject           string // https://api.commerce7.com/v1/tag-x-object/{:obj} - Requires an object appended; either "customer" or "order"
	LoyaltyTransaction   string // "https://api.commerce7.com/v1/loyalty-transaction" - requires Commerce7 Loyalty Extension https://documentation.commerce7.com/loyalty-feature-overview
	Vendor               string // "https://api.commerce7.com/v1/vendor" - example: https://api.commerce7.com/v1/vendor?q=myvendor
	Collections          string // "https://api.commerce7.com/v1/collection" - example: https://api.commerce7.com/v1/collection?type=Manual
	Department           string // "https://api.commerce7.com/v1/department" - example: https://api.commerce7.com/v1/department?q=wine
	Setting              string // "https://api.commerce7.com/v1/setting" - Used for grabbing the winery info/settings from a given tenant.
	WineAppellation      string // "https://api.commerce7.com/v1/wine-appellation" - Get valid wine appellations.
	WineVarietal         string // "https://api.commerce7.com/v1/wine-varietal" - Get valid wine varietals.
}

func GetEndpoints(baseURL string) *endpoints {
	return &endpoints{
		Auth:                 baseURL + "/account/user",
		Cart:                 baseURL + "/cart",
		Customer:             baseURL + "/customer",
		CustomerId:           baseURL + "/customer/{:id}",
		ClubMembership:       baseURL + "/club-membership",
		FulfillmentAll:       baseURL + "/order/{:id}/fulfillment/all",
		Fulfillment:          baseURL + "/order/{:id}/fulfillment",
		GiftCard:             baseURL + "/gift-card",
		GiftCardTransaction:  baseURL + "/gift-card-transaction",
		Inventory:            baseURL + "/inventory",
		InventoryTransaction: baseURL + "/inventory-transaction",
		MetaDataConfig:       baseURL + "/meta-data-config/",
		MetaDataConfigObj:    baseURL + "/meta-data-config/{:obj}",
		Order:                baseURL + "/order",
		Product:              baseURL + "/product",
		Tag:                  baseURL + "/tag",
		TagXObject:           baseURL + "/tag-x-object/{:obj}",
		LoyaltyTransaction:   baseURL + "/loyalty-transaction",
		Vendor:               baseURL + "/vendor",
		Collections:          baseURL + "/collection",
		Department:           baseURL + "/department",
		Setting:              baseURL + "/setting",
		WineAppellation:      baseURL + "/wine-appellation",
		WineVarietal:         baseURL + "/wine-varietal",
	}
}

var Endpoints = GetEndpoints(API_URL)
var EndpointsV2 = GetEndpoints(API_URL_V2)
