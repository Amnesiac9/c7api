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

// Product Type
const (
	ProductTypeGeneralMerchandise = "General Merchandise"
	ProductTypeTasting            = "Tasting"
	ProductTypeWine               = "Wine"
	ProductTypeCannabis           = "Cannabis"
	ProductTypeBundle             = "Bundle"
	ProductTypeReservation        = "Reservation"
	ProductTypeEventTicket        = "Event Ticket"
	ProductTypeGiftCard           = "Gift Card"
	ProductTypeCollateral         = "Collateral"
	ProductTypeRebate             = "Rebate"
)

// Webhook Action
const (
	WebhookActionCreate     = "Create"
	WebhookActionUpdate     = "Update"
	WebhookActionBulkUpdate = "Bulk Update"
	WebhookActionDelete     = "Delete"
)

// Club Membership Status
const (
	ClubMembershipStatusActive    = "Active"
	ClubMembershipStatusCancelled = "Cancelled"
)

// GiftCard Status
const (
	GiftCardStatusActive    = "Active"
	GiftCardStatusCancelled = "Cancelled"
)

// GiftCard Type

const (
	GiftCardTypeVirtual  = "Virtual"
	GiftCardTypePhysical = "Physical"
)
