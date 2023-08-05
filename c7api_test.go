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

	jsonBytes2, err := GetJsonFromC7(&urlString, &tenant, &badAuth)
	if err == nil {
		t.Error("Error, got JSON from C7 with bad auth: ", err.Error())
		return
	}

	if jsonBytes2 != nil {
		t.Error("JSON Bytes should be nil with bad auth")
		return
	}

}

func TestPostJsonToC7(t *testing.T) {

	urlString := "https://api.commerce7.com/v1/order"
	urlStringFulfillment := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/all"
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

	//TODO: Maybe have post return the response or status code
	jsonBytes, statusCode, err := PostJsonToC7(&urlString, &tenant, &blankBytes, &goodAuth)
	if err != nil && statusCode == 0 {
		t.Error("Error getting JSON from C7: ", err.Error())
		return
	}

	if statusCode != 422 {
		t.Error("Status code should be 422 with bad JSON, got: ", statusCode)
		return
	}

	if jsonBytes == nil {
		t.Error("jsonBytes2 from C7 should not be nil")
		return
	}

	jsonBytes2, statusCode, err := PostJsonToC7(&urlString, &tenant, &blankBytes, &badAuth)
	if err != nil {
		t.Error("Error getting JSON from C7: ", err.Error())
		return
	}

	if statusCode != 401 {
		t.Error("Status code should be 401 with bad auth, got: ", statusCode)
		return
	}

	if jsonBytes2 == nil {
		t.Error("jsonBytes2 from C7 should not be nil")
		return
	}

	jsonBytes3, statusCode, err := PostJsonToC7(&urlStringFulfillment, &tenant, &goodBytes, &goodAuth)
	if err != nil {
		t.Error("Got Error with good auth and good bytes: ", err.Error())
		if jsonBytes3 != nil {
			fmt.Println(string(*jsonBytes3))
		}
		return
	}

	if statusCode != 422 {
		t.Error("Status code should be 422 with good auth and good bytes, got: ", statusCode)
		return
	}

	if jsonBytes3 == nil {
		t.Error("jsonBytes3 from C7 should not be nil")
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
