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

// Struct to return as response to C7 Order Details Page
//
//	{
//		   "icon": url of icon on their server,
//		   "subTitle": "",
//		   "footer": "",
//		   "title": "",
//		   "variant": null [null, 'success', 'warning', 'error']
//	   }
type OrderDetailsStatusCard struct {
	Icon     string `json:"icon"`
	SubTitle string `json:"subTitle"`
	Footer   string `json:"footer"`
	Title    string `json:"title"`
	Variant  string `json:"variant,omitempty"` // Can be null, "success", "warning", or "error"
}
