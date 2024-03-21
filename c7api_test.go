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
	AppAuthEncoded = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		os.Getenv("appid"),
		os.Getenv("appkey"),
	)))
	testTenant = os.Getenv("testTenant")
)

func TestGetC7_New(t *testing.T) {

	urlStringOrders := "https://api.commerce7.com/v1/order?orderPaidDate=btw:2023-07-29T07:00:00.000Z|2023-07-31T06:59:59.999Z"
	tenant := testTenant
	goodAuth := AppAuthEncoded
	//badAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth"))

	type testCase struct {
		name          string
		method        string
		url           string
		body          []byte
		tenant        string
		auth          string
		attempts      int
		expectedCode  int
		expectedBytes []byte
	}

	testCases := []testCase{
		{
			name:          "Good GET",
			method:        "GET",
			url:           urlStringOrders,
			body:          nil,
			tenant:        tenant,
			auth:          goodAuth,
			attempts:      0,
			expectedCode:  200,
			expectedBytes: nil, // TODO: add expected bytes
		},
		{
			name:          "Bad Auth GET",
			method:        "GET",
			url:           urlStringOrders,
			body:          nil,
			tenant:        tenant,
			auth:          "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth")),
			attempts:      0,
			expectedCode:  401,
			expectedBytes: nil,
		},
	}

	for _, tc := range testCases {

		t.Log("Test", tc.name)

		jsonBytes, err := NewRequest(tc.method, &tc.url, &tc.body, tc.tenant, tc.auth, tc.attempts)
		if err != nil && err.(C7Error).StatusCode != tc.expectedCode {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, " Expected status code: ", tc.expectedCode, " got: ", err.(C7Error).StatusCode)
		}

		if tc.expectedBytes != nil && string(*jsonBytes) != string(tc.expectedBytes) {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, "Expected: ", string(tc.expectedBytes), " got: ", string(*jsonBytes))
			return
		}

	}

}

func TestPostC7_New(t *testing.T) {

	urlStringFulfillment := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/all"
	orderNumber := 1235
	orderId := "034e6096-429d-452c-b258-5d37a1522934"
	tenant := testTenant
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

	testCases := []struct {
		name          string
		method        string
		url           string
		body          []byte
		tenant        string
		auth          string
		attempts      int
		expectedCode  int
		expectedBytes []byte
	}{
		{"Good Post", "POST", urlStringFulfillment, goodBytes, tenant, goodAuth, 0, 200, []byte(`{"id":"034e6096-429d-452c-b258-5d37a1522934","orderSubmittedDate":"2023-07-30T20:44:32.725Z","orderPaidDate":"2023-07-30T20:44:32.725Z","orderFulfilledDate":"2023-07-30T10:59:32.000Z","orderNumber":1235,`)},
		{"Bad Auth POST", "POST", urlStringFulfillment, nil, tenant, badAuth, 0, 401, []byte(`{"statusCode":401,"type":"unauthorized","message":"Unauthenticated User","errors":[]}`)},
		{"Already Fulfileld POST", "POST", urlStringFulfillment, goodBytes, tenant, goodAuth, 0, 422, []byte(`{"statusCode":422,"type":"validationError","message":"Can not fulfill an order that is marked Fulfilled"}`)},
		{"Blank POST", "POST", urlStringFulfillment, blankBytes, tenant, goodAuth, 0, 422, []byte(`{"statusCode":422,"type":"validationError","message":"One or more elements is missing or invalid","errors":[{"field":"type","message":"required"},{"field":"fulfillmentDate","message":"required"},{"field":"packageCount","message":"required"}]}`)},
	}

	// Delete previous fulfillment for test
	fulfillmentId, err := GetFulfillmentId(orderNumber, testTenant, AppAuthEncoded, 1)
	if err != nil {
		t.Error("Error getting fulfillment id: ", err.Error())
		return
	}

	t.Log("Deleting Fulfillment ID: ", fulfillmentId)

	_, err = DeleteC7Fulfillment(orderId, fulfillmentId, testTenant, AppAuthEncoded, 1)
	if err != nil {
		t.Error("Error deleting fulfillment: ", err.Error())
		return
	}

	for _, tc := range testCases {

		t.Log("Test", tc.name)

		jsonBytes, err := NewRequest(tc.method, &tc.url, &tc.body, tc.tenant, tc.auth, tc.attempts)
		if err != nil && err.(C7Error).StatusCode != tc.expectedCode {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, " Expected status code: ", tc.expectedCode, " got: ", err.(C7Error).StatusCode)
		}

		if tc.expectedBytes != nil && string(*jsonBytes)[:50] != string(tc.expectedBytes)[:50] {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, "Expected: ", string(tc.expectedBytes)[:50], " got: ", string(*jsonBytes)[:50])
			return
		}

	}

}

func TestDeleteC7_New(t *testing.T) {

	//urlStringDelete := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/{fulfillmentId}"
	urlStringFulfillment := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/all"
	tenant := testTenant
	goodAuth := AppAuthEncoded
	badAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth"))
	orderNumber := 1235

	//blankBytes := []byte("")
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

	fulfillmentId, err := GetFulfillmentId(orderNumber, testTenant, AppAuthEncoded, 1)
	if err != nil {
		t.Error("Error getting fulfillment id: ", err.Error())
		return
	}

	t.Log("Deleting Fulfillment ID: ", fulfillmentId)

	urlStringDelete := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/" + fulfillmentId

	testCases := []struct {
		name          string
		method        string
		url           string
		body          []byte
		tenant        string
		auth          string
		attempts      int
		expectedCode  int
		expectedBytes []byte
	}{
		{"Good DELETE", "DELETE", urlStringDelete, nil, tenant, goodAuth, 0, 200, []byte(`{"id":"034e6096-429d-452c-b258-5d37a1522934","orderSubmittedDate":"2023-07-30T20:44:32.725Z","orderPaidDate":"2023-07-30T20:44:32.725Z","orderFulfilledDate":null,"orderNumber":1235,`)},
		{"Bad Auth DELETE", "DELETE", urlStringDelete, nil, tenant, badAuth, 0, 401, []byte(`{"statusCode":401,"type":"unauthorized","message":"Unauthenticated User","errors":[]}`)},
		{"No Fulfillment to DELETE", "DELETE", urlStringDelete, nil, tenant, goodAuth, 0, 422, []byte(`{"statusCode":422,"type":"processingError","message":"Fulfillment not found"}`)},
	}

	for _, tc := range testCases {

		t.Log("Test", tc.name)

		jsonBytes, err := NewRequest(tc.method, &tc.url, &tc.body, tc.tenant, tc.auth, tc.attempts)
		if err != nil && err.(C7Error).StatusCode != tc.expectedCode {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, " Expected status code: ", tc.expectedCode, " got: ", err.(C7Error).StatusCode)
		}

		if tc.expectedBytes != nil && string(*jsonBytes)[:50] != string(tc.expectedBytes)[:50] {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, "Expected: ", string(tc.expectedBytes)[:50], " got: ", string(*jsonBytes)[:50])
			break
		}

	}

	t.Log("Adding Fulfillment for test TestDeleteC7_New")

	// Post previous fulfillment for test
	jsonBytes, err := NewRequest("POST", &urlStringFulfillment, &goodBytes, tenant, goodAuth, 1)
	if err != nil || jsonBytes == nil {
		t.Error("Error posting fulfillment: ", err.Error())
		return
	}

}

/// OLD TESTS ///

func TestGetJSONFromC7(t *testing.T) {

	urlString := "https://api.commerce7.com/v1/order?orderPaidDate=btw:2023-07-29T07:00:00.000Z|2023-07-31T06:59:59.999Z"
	tenant := testTenant
	goodAuth := AppAuthEncoded
	badAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth"))

	jsonBytes, err := GetJsonFromC7(&urlString, tenant, goodAuth, 0)
	if err != nil {
		t.Error("Error getting JSON from C7: ", err.Error())
		return
	}

	if jsonBytes == nil {
		t.Error("JSON from C7 is empty")
		return
	}

	jsonBytes, err = GetJsonFromC7(&urlString, tenant, badAuth, 3)
	if err == nil {
		t.Error("Error, did not get err with bad auth.")
		return
	}
	fmt.Println(err.Error())

	if err.Error() != `status code: 401, error: {"statusCode":401,"type":"unauthorized","message":"Unauthenticated User","errors":[]}` {
		t.Error("Error, expected: ", `Status Code: 401, C7 Error: {"statusCode":401,"type":"unauthorized","message":"Unauthenticated User","errors":[]}`, " got: ", err.Error())
		return
	}

	if err.(C7Error).StatusCode != 401 {
		t.Error("Error, expected status code 401, got: ", err.(C7Error).StatusCode)
		return
	}

	if jsonBytes == nil {
		t.Error("JSON Bytes should not be nil with bad auth")
		return
	}

	_, err = GetJsonFromC7(nil, "", "", 0)
	if err == nil {
		t.Error("Error should not be nil with nil params.")
		return
	}

}

func TestPostJsonToC7(t *testing.T) {

	//urlString := "https://api.commerce7.com/v1/order"
	urlStringFulfillment := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/all"
	//urlStringRemoveFulfillment := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/f1439243-ceee-4ed6-b08a-4bd12f36c63e"
	orderNumber := 1235
	orderId := "034e6096-429d-452c-b258-5d37a1522934"
	tenant := testTenant
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

	testCases := []struct {
		urlString     string
		tenant        string
		bytes         []byte
		auth          string
		attempts      int
		expectedCode  int
		expectedBytes []byte
	}{
		// 1 Blank Bytes
		{urlStringFulfillment, tenant, blankBytes, goodAuth, 2, 422, []byte(`{"statusCode":422,"type":"validationError","message":"One or more elements is missing or invalid","errors":[{"field":"type","message":"required"},{"field":"fulfillmentDate","message":"required"},{"field":"packageCount","message":"required"}]}`)},
		// 2 Bad Auth
		{urlStringFulfillment, tenant, goodBytes, badAuth, 2, 401, []byte(`{"statusCode":401,"type":"unauthorized","message":"Unauthenticated User","errors":[]}`)},
		// 3 Good Auth, Good Bytes
		{urlStringFulfillment, tenant, goodBytes, goodAuth, 2, 200, []byte(`{"id":"034e6096-429d-452c-b258-5d37a1522934","orderSubmittedDate":"2023-07-30T20:44:32.725Z","orderPaidDate":"2023-07-30T20:44:32.725Z","orderFulfilledDate":"2023-07-30T10:59:32.000Z","orderNumber":1235,`)},
		// 4 Good Auth, Good Bytes, already fulfilled
		{urlStringFulfillment, tenant, goodBytes, goodAuth, 0, 422, []byte(`{"statusCode":422,"type":"validationError","message":"Can not fulfill an order that is marked Fulfilled"}`)},
	}

	// Delete previous fulfillment for test
	fulfillmentId, err := GetFulfillmentId(orderNumber, testTenant, AppAuthEncoded, 1)
	if err != nil {
		t.Error("Error getting fulfillment id: ", err.Error())
		return
	}

	t.Log("Fulfillment ID: ", fulfillmentId)

	_, err = DeleteC7Fulfillment(orderId, fulfillmentId, testTenant, AppAuthEncoded, 1)
	if err != nil {
		t.Error("Error deleting fulfillment: ", err.Error())
		return
	}

	for i, testCase := range testCases {
		jsonBytes, err := PostJsonToC7(&testCase.urlString, &testCase.bytes, testCase.tenant, testCase.auth, testCase.attempts)
		if err != nil && err.(C7Error).StatusCode != testCase.expectedCode {
			t.Error("TestPostJsonToC7, test case: ", i+1, " Expected status code: ", testCase.expectedCode, " got: ", err.(C7Error).StatusCode)
		}
		if string(*jsonBytes)[:50] != string(testCase.expectedBytes)[:50] {
			t.Error("TestPostJsonToC7, test case: ", i+1, "Expected: ", string(testCase.expectedBytes)[:50], " got: ", string(*jsonBytes)[:50])
			return
		}
	}

}

func TestDeleteFromC7(t *testing.T) {

	urlString := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/f1439243-ceee-4ed6-b08a-4bd12f36c63e"
	tenant := testTenant
	goodAuth := AppAuthEncoded
	badAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth"))

	// Test 1
	bytes, err := DeleteFromC7(&urlString, tenant, badAuth, 2)
	if err == nil {
		t.Error("Error should not be nil with bad auth.")
		return
	}

	if err.(C7Error).StatusCode != 401 {
		t.Error("Status code should be 401 with bad auth, got: ", err.(C7Error).StatusCode)
		return
	}

	if bytes == nil {
		t.Error("bytes from C7 should not be nil with bad auth")
		return
	}

	// Test 2
	bytes2, err := DeleteFromC7(&urlString, tenant, goodAuth, 0)
	if err == nil {
		t.Error("Error should not be nil with good auth.")
		return
	}

	if err.(C7Error).StatusCode != 422 {
		t.Error("Status code should be 422 with good auth, got: ", err.(C7Error).StatusCode)
		return
	}

	if bytes2 == nil {
		t.Error("bytes2 from C7 should not be nil")
		return
	}

	_, err = DeleteFromC7(nil, "", "", 0)
	if err == nil {
		t.Error("Error should not be nil with nil params.")
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

func TestGetFulfillmentId(t *testing.T) {
	// testParams := []int{1232, 1005, 999, 1239}
	// expected := []string{"9475723a-8f11-4111-9234-852d85813581", "139826d7-348d-4a3b-aca0-5466e7462e79", "", ""}

	testCases := []struct {
		orderNumber  int
		auth         string
		expectedCode int
		expectedId   string
	}{
		{1232, AppAuthEncoded, 0, "9475723a-8f11-4111-9234-852d85813581"},
		{1007, AppAuthEncoded, 0, "030c8ab2-354b-41a7-acd4-a13bedc70dd5"},
		{999, AppAuthEncoded, 0, ""},
		{1239, AppAuthEncoded, 0, ""},
		{1232, "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth")), 401, ""},
	}

	// for i, param := range testParams {
	// 	result, _ := GetFulfillmentId(param, testTenant, AppAuthEncoded)
	// 	if result != expected[i] {
	// 		t.Error("Expected ", expected[i], " got ", result)
	// 	}
	// }

	for i, testCase := range testCases {
		resultId, err := GetFulfillmentId(testCase.orderNumber, testTenant, testCase.auth, 1)
		if err, ok := err.(C7Error); ok {
			if err.StatusCode != testCase.expectedCode {
				t.Error("Test case:", i+1, "expected status code:", testCase.expectedCode, "got:", err.StatusCode)
			}
		}
		if resultId != testCase.expectedId {
			t.Error("Test case:", i+1, "expected id:", testCase.expectedId, "got:", resultId)
		}
	}
}

func Test_IsCarrierSupported(t *testing.T) {
	testCases := []struct {
		carrier  string
		expected bool
	}{
		{"UPS", true},
		{"FedEx", true},
		{"GSO", true},
		{"ATS Healthcare", true},
		{"Australia Post", true},
		{"USPS", false},
		{"Canada Post", false},
		{"DHL", false},
		{"", false},
	}

	for i, testCase := range testCases {
		result := IsCarrierSupported(testCase.carrier)
		if result != testCase.expected {
			t.Error("Test case: ", i+1, " expected: ", testCase.expected, " got: ", result)
		}
	}
}
