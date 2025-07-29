package c7api

// Order Status
const (
	OrderFulfillmentStatusNotFulfilled          = "Not Fulfilled"
	OrderFulfillmentStatusNoFulfillmentRequired = "No Fulfillment Required"
	OrderFulfillmentStatusPartiallyFulfilled    = "Partially Fulfilled"
	OrderFulfillmentStatusFulfilled             = "Fulfilled"
)

// Delivery Method
const (
	OrderDeliveryMethodShip     = "Ship"
	OrderDeliveryMethodPickup   = "Pickup"
	OrderDeliveryMethodCarryout = "Carry Out"
)

// Channel

const (
	OrderChannelPOS     = "POS"
	OrderChannelInbound = "Inbound"
	OrderChannelClub    = "Club"
	OrderChannelWeb     = "Web"
)
