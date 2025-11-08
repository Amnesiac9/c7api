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
	ClubMembershipStatusOnHold    = "On Hold"
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

// MetaDataConfig Object Types
const (
	MetaDataConfigObjectAllocation      = "allocation"
	MetaDataConfigObjectClubMembership  = "club-membership"
	MetaDataConfigObjectCollection      = "collection"
	MetaDataConfigObjectCustomer        = "customer"
	MetaDataConfigObjectCustomerAddress = "customer-address"
	MetaDataConfigObjectOrder           = "order"
	MetaDataConfigObjectProduct         = "product"
	MetaDataConfigObjectReservation     = "reservation"
	MetaDataConfigObjectExperience      = "experience"
)

// Tag X Object Types
const (
	TagXObjectTypeCustomer = "customer"
	TagXObjectTypeOrder    = "order"
)

func IsValidTagXObjectType(objectType string) bool {
	switch objectType {
	case TagXObjectTypeCustomer,
		TagXObjectTypeOrder:
		return true
	default:
		return false
	}
}

// Variant Tax Types
const (
	TaxTypeWine               = "Wine"
	TaxTypeGeneralMerchandise = "General Merchandise"
	TaxTypeFood               = "Food"
	TaxTypePreparedFood       = "Prepared Food"
	TaxTypeBooks              = "Books"
	TaxTypeNotTaxable         = "Not Taxable"
)

func IsValidTaxType(taxType string) bool {
	switch taxType {
	case TaxTypeWine,
		TaxTypeGeneralMerchandise,
		TaxTypeFood,
		TaxTypePreparedFood,
		TaxTypeBooks,
		TaxTypeNotTaxable:
		return true
	default:
		return false
	}
}

const (
	AdminStatusAvailable    = "Available"
	AdminStatusNotAvailable = "Not Available"
)

func IsValidAdminStatus(status string) bool {
	switch status {
	case AdminStatusAvailable,
		AdminStatusNotAvailable:
		return true
	default:
		return false
	}
}

const (
	WebStatusAvailable    = "Available"
	WebStatusNotAvailable = "Not Available"
	WebStatusRetired      = "Retired"
)

func IsValidWebStatus(status string) bool {
	switch status {
	case WebStatusAvailable,
		WebStatusNotAvailable,
		WebStatusRetired:
		return true
	default:
		return false
	}
}
