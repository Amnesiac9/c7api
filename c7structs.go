package c7api

import (
	"time"
)

type FulfillmentAllItems struct {
	SendTransactionEmail bool      `json:"sendTransactionEmail"`
	Type                 string    `json:"type"`
	FulfillmentDate      time.Time `json:"fulfillmentDate"`
	Shipped              *struct {
		TrackingNumbers []string `json:"trackingNumbers"`
		Carrier         string   `json:"carrier"`
	} `json:"shipped"`
	PackageCount int `json:"packageCount"`
}

type NewUser struct {
	TenantID   string `json:"tenantId"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	TimeZone   string `json:"timeZone"`
	WeightUnit string `json:"weightUnit"`
	User       struct {
		UserID    string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	} `json:"user"`
}
