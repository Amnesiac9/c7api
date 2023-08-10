package c7api

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")
var (
	// App C7 Basic Auth Credentials
	AppAuth = fmt.Sprintf("%s:%s",
		os.Getenv("appid"),
		os.Getenv("appkey"),
	)

	AppAuthEncoded = "Basic " + base64.StdEncoding.EncodeToString([]byte(AppAuth))
)

func TestGetJSONFromC7(t *testing.T) {

	urlString := "https://api.commerce7.com/v1/order?orderPaidDate=btw:2023-07-29T07:00:00.000Z|2023-07-31T06:59:59.999Z"
	tenant := "egyptian-thread-company"
	goodAuth := AppAuthEncoded
	badAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth"))

	jsonBytes, err := GetJsonFromC7(&urlString, &tenant, &goodAuth)
	if err != nil {
		t.Error("Error getting JSON from C7: ", err.Error())
		return
	}

	if jsonBytes == nil {
		t.Error("JSON from C7 is empty")
		return
	}

	jsonBytes, err = GetJsonFromC7(&urlString, &tenant, &badAuth)
	if err == nil {
		t.Error("Error, did not get err with bad auth: ", err.Error())
		return
	}
	fmt.Println(err.Error())

	if err.Error() != `Status Code: 401, C7 Error: {"statusCode":401,"type":"unauthorized","message":"Unauthenticated User","errors":[]}` {
		t.Error("Error, expected: ", `Status Code: 401, C7 Error: {"statusCode":401,"type":"unauthorized","message":"Unauthenticated User","errors":[]}`, " got: ", err.Error())
		return
	}

	if err.(C7Error).StatusCode != 401 {
		t.Error("Error, expected status code 401, got: ", err.(C7Error).StatusCode)
		return
	}

	if jsonBytes != nil {
		t.Error("JSON Bytes should be nil with bad auth")
		return
	}

}

func TestPostJsonToC7(t *testing.T) {

	urlString := "https://api.commerce7.com/v1/order"
	urlStringFulfillment := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/all"
	//urlStringRemoveFulfillment := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/f1439243-ceee-4ed6-b08a-4bd12f36c63e"
	tenant := "egyptian-thread-company"
	goodAuth := AppAuthEncoded
	badAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth"))

	blankBytes := []byte("")
	goodBytes := []byte(`{
		"sendTransactionEmail": false,
		"type": "Shipped",
		"fulfillmentDate": "2023-07-30T10:59:32.000Z",
		"shipped": {
			"trackingNumbers": ["1Z525R5EA803600000"],
			"carrier": "UPS"
		},
		"packageCount": 1
	}`)

	jsonBytes, err := PostJsonToC7(&urlString, &tenant, &blankBytes, &goodAuth)
	if err == nil {
		t.Error("Error should not be nil with blank bytes.")
		return
	}

	if err.(C7Error).StatusCode != 422 {
		t.Error("Status code should be 422 with blank JSON, got: ", err.(C7Error).StatusCode)
		return
	}

	if jsonBytes == nil {
		t.Error("jsonBytes from C7 should not be nil")
		return
	}

	if err.(C7Error).Err == nil {
		t.Error("C7Error.Error should not be nil")
		return
	}

	// Test with bad auth
	jsonBytes2, err := PostJsonToC7(&urlString, &tenant, &goodBytes, &badAuth)
	if err == nil {
		t.Error("Error should not be nil with bad auth.")
		return
	}

	if err.(C7Error).StatusCode != 401 {
		t.Error("Status code should be 401 with bad auth, got: ", err.(C7Error).StatusCode)
		return
	}

	if jsonBytes2 == nil {
		t.Error("jsonBytes2 from C7 should not be nil with bad auth")
		return
	}

	jsonBytes3, err := PostJsonToC7(&urlStringFulfillment, &tenant, &goodBytes, &goodAuth)
	if err == nil {
		t.Error("Error should not be nil with good auth and good bytes.")
		return
	}
	// if err != nil {
	// 	t.Error("Got error with good auth and good bytes: ", err.Error())
	// 	if jsonBytes3 != nil {
	// 		fmt.Println(string(*jsonBytes3))
	// 	}
	// 	return
	// }

	if err.(C7Error).StatusCode != 422 {
		t.Error("Status code should be 422 with good auth and good bytes, got: ", err.(C7Error).StatusCode)
		return
	}

	if jsonBytes3 == nil {
		t.Error("jsonBytes3 from C7 should not be nil")
		return
	}

	expected := `{"statusCode":422,"type":"validationError","message":"Can not fulfill an order that is marked Fulfilled"}`
	if string(*jsonBytes3) != expected {
		t.Error("Expected: ", expected, " got: ", string(*jsonBytes3))
		return
	}

}

func TestFormatDatesForC7(t *testing.T) {
	testParams := []string{"01/02/2006 15:04", "01/02/2025 15:04", "01/02/2006", ""}
	expected := []string{"2006-01-02T15:04:00.000Z", "2025-01-02T15:04:00.000Z", "01/02/2006", ""}

	for i, param := range testParams {
		result, err := FormatDatesForC7(param)
		if result != expected[i] {
			t.Error("Expected ", expected[i], " got ", result)
		}
		if err != nil {
			if err.Error() != `parsing time "01/02/2006" as "01/02/2006 15:04": cannot parse "" as "15"` && err.Error() != `date is empty` {
				t.Error("Expected error, got ", err.Error())
			}
		}
	}
}
