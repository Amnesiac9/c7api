package c7api

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

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

type rateLimiterMock struct{}

func (rl *rateLimiterMock) Wait() {
	time.Sleep(1 * time.Millisecond)
}

func TestGetC7_Request(t *testing.T) {

	urlStringOrders := "https://api.commerce7.com/v1/order?orderPaidDate=btw:2023-07-29T07:00:00.000Z|2023-07-31T06:59:59.999Z"
	tenant := testTenant
	goodAuth := AppAuthEncoded
	//badAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth"))

	type testCase struct {
		name               string
		method             string
		url                string
		body               []byte
		tenant             string
		auth               string
		attempts           int
		expectedStatusCode int
		expectedErrorText  string
	}

	testCases := []testCase{
		{
			name:               "Good GET",
			method:             "GET",
			url:                urlStringOrders,
			body:               nil,
			tenant:             tenant,
			auth:               goodAuth,
			attempts:           0,
			expectedStatusCode: 200,
		},
		{
			name:               "Bad Auth GET",
			method:             "GET",
			url:                urlStringOrders,
			body:               nil,
			tenant:             tenant,
			auth:               "Basic " + base64.StdEncoding.EncodeToString([]byte("bad:auth")),
			attempts:           0,
			expectedStatusCode: 401,
			expectedErrorText:  "reponse status not within 200-299: 401 Unauthorized",
		},
	}

	for _, tc := range testCases {

		//t.Log("Test", tc.name)

		resp, err := Request(tc.method, tc.url, &tc.body, tc.tenant, tc.auth, true)
		if err != nil && err.Error() != tc.expectedErrorText {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, " unexpected error returned: ", err.Error())
			return
		}

		if resp == nil {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, " response is nil")
			return
		}

		if resp.StatusCode != tc.expectedStatusCode {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, " Expected status code: ", tc.expectedStatusCode, " got: ", resp.StatusCode)
		}

	}

}

func TestGetC7_New(t *testing.T) {

	urlStringOrders := "https://api.commerce7.com/v1/order"
	queries := map[string]string{
		"orderPaidDate": "btw:2023-07-29T07:00:00.000Z|2023-07-31T06:59:59.999Z",
	}
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
		rl            genericRateLimiter
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
			rl:            &rateLimiterMock{},
			expectedCode:  401,
			expectedBytes: nil,
		},
	}

	for _, tc := range testCases {

		//t.Log("Test", tc.name)

		jsonBytes, err := RequestWithRetryAndRead(tc.method, tc.url, queries, &tc.body, tc.tenant, tc.auth, tc.attempts, tc.rl)
		if err != nil {
			if c7err, ok := err.(*C7Error); ok {
				if c7err.StatusCode != tc.expectedCode {
					t.Error("TestGetJSONFromC7, test case: ", tc.name, " Expected status code: ", tc.expectedCode, " got: ", c7err.StatusCode)
					return
				}
			} else {
				t.Error("TestGetJSONFromC7, test case: ", tc.name, " Expected status code: ", tc.expectedCode, " got: ", err.Error())
				return
			}
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
		name               string
		method             string
		url                string
		body               []byte
		tenant             string
		auth               string
		attempts           int
		rl                 genericRateLimiter
		expectedStatusCode int
		expectedBytes      []byte
	}{
		{"Good Post", "POST", urlStringFulfillment, goodBytes, tenant, goodAuth, 0, &rateLimiterMock{}, 200, []byte(`{"id":"034e6096-429d-452c-b258-5d37a1522934","orderSubmittedDate":"2023-07-30T20:44:32.725Z","orderPaidDate":"2023-07-30T20:44:32.725Z","orderFulfilledDate":"2023-07-30T10:59:32.000Z","orderNumber":1235,`)},
		{"Bad Auth POST", "POST", urlStringFulfillment, nil, tenant, badAuth, 0, nil, 401, []byte(`{"statusCode":401,"type":"unauthorized","message":"Unauthenticated User","errors":[]}`)},
		{"Already Fulfileld POST", "POST", urlStringFulfillment, goodBytes, tenant, goodAuth, 0, nil, 422, []byte(`{"statusCode":422,"type":"validationError","message":"Can not fulfill an order that is marked Fulfilled"}`)},
		{"Blank POST", "POST", urlStringFulfillment, blankBytes, tenant, goodAuth, 0, nil, 422, []byte(`{"statusCode":422,"type":"validationError","message":"One or more elements is missing or invalid","errors":[{"field":"type","message":"required"},{"field":"fulfillmentDate","message":"required"},{"field":"packageCount","message":"required"}]}`)},
	}

	// Delete previous fulfillment for test
	fulfillmentIds, err := GetFulfillmentIds(orderNumber, testTenant, AppAuthEncoded, 1, nil)
	if err != nil {
		t.Error("Error getting fulfillment id: ", err.Error())
		return
	}

	for _, id := range fulfillmentIds {
		//t.Log("Deleting Fulfillment ID: ", id)

		_, err = DeleteFulfillmentById(orderId, id, testTenant, AppAuthEncoded, 1, nil)
		if err != nil {
			t.Error("Error deleting fulfillment: ", err.Error())
			return
		}
	}

	for _, tc := range testCases {

		t.Log("Test", tc.name)

		jsonBytes, err := RequestWithRetryAndRead(tc.method, tc.url, nil, &tc.body, tc.tenant, tc.auth, tc.attempts, tc.rl)
		if err != nil && err.(*C7Error).StatusCode != tc.expectedStatusCode {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, " Expected status code: ", tc.expectedStatusCode, " got: ", err.(*C7Error).StatusCode)
		}

		// TODO: Unmarshal these and compare that way to get exact matches.
		if tc.expectedBytes != nil && string(*jsonBytes)[:50] != string(tc.expectedBytes)[:50] {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, "Expected: ", string(tc.expectedBytes)[:50], " got: ", string(*jsonBytes)[:50])
			return
		}

	}

}

func TestDeleteC7_New(t *testing.T) {

	//urlStringDelete := "https://api.commerce7.com/v1/order/034e6096-429d-452c-b258-5d37a1522934/fulfillment/{fulfillmentId}"
	urlStringFulfillment := Endpoints.Order + "/034e6096-429d-452c-b258-5d37a1522934/fulfillment/all"
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

	fulfillmentIds, err := GetFulfillmentIds(orderNumber, testTenant, AppAuthEncoded, 1, nil)
	if err != nil {
		t.Error("Error getting fulfillment id: ", err.Error())
		return
	}

	if len(fulfillmentIds) != 1 {
		t.Errorf("expected fulfillmentIds length of 1, got: %d", len(fulfillmentIds))
		return
	}

	fulfillmentId := fulfillmentIds[0]

	t.Log("Deleting Fulfillment ID: ", fulfillmentId)

	urlStringDelete := Endpoints.Order + "/034e6096-429d-452c-b258-5d37a1522934/fulfillment/" + fulfillmentId

	testCases := []struct {
		name          string
		method        string
		url           string
		body          []byte
		tenant        string
		auth          string
		attempts      int
		rl            genericRateLimiter
		expectedCode  int
		expectedBytes []byte
	}{
		{"Good DELETE", "DELETE", urlStringDelete, nil, tenant, goodAuth, 0, &rateLimiterMock{}, 200, []byte(`{"id":"034e6096-429d-452c-b258-5d37a1522934","orderSubmittedDate":"2023-07-30T20:44:32.725Z","orderPaidDate":"2023-07-30T20:44:32.725Z","orderFulfilledDate":null,"orderNumber":1235,`)},
		{"Bad Auth DELETE", "DELETE", urlStringDelete, nil, tenant, badAuth, 0, nil, 401, []byte(`{"statusCode":401,"type":"unauthorized","message":"Unauthenticated User","errors":[]}`)},
		{"No Fulfillment to DELETE", "DELETE", urlStringDelete, nil, tenant, goodAuth, 0, nil, 422, []byte(`{"statusCode":422,"type":"processingError","message":"Fulfillment not found"}`)},
	}

	for _, tc := range testCases {

		//t.Log("Test", tc.name)

		jsonBytes, err := RequestWithRetryAndRead(tc.method, tc.url, nil, &tc.body, tc.tenant, tc.auth, tc.attempts, tc.rl)
		if err != nil && err.(*C7Error).StatusCode != tc.expectedCode {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, " Expected status code: ", tc.expectedCode, " got: ", err.(*C7Error).StatusCode)
		}

		if tc.expectedBytes != nil && string(*jsonBytes)[:50] != string(tc.expectedBytes)[:50] {
			t.Error("TestGetJSONFromC7, test case: ", tc.name, "Expected: ", string(tc.expectedBytes)[:50], " got: ", string(*jsonBytes)[:50])
			break
		}

	}

	t.Log("Re-Adding Fulfillment for test TestDeleteC7_New")

	// Post previous fulfillment for test
	jsonBytes, err := RequestWithRetryAndRead("POST", urlStringFulfillment, nil, &goodBytes, tenant, goodAuth, 1, nil)
	if err != nil || jsonBytes == nil {
		t.Error("Error posting fulfillment: ", err.Error())
		return
	}

}

func Test_GetFulfillments(t *testing.T) {

	orderNumber := 1235

	fulfillments, err := GetFulfillmentsByOrderNumber(orderNumber, testTenant, AppAuthEncoded, 1, nil)
	if err != nil {
		t.Error("error getting fulfillments: ", err.Error())
		return
	}

	if fulfillments == nil {
		t.Error("error: fulfillments are nil")
		return
	}

	if len(*fulfillments) != 1 {
		t.Errorf("expected fulfillmentIds length of 1, got: %d", len(*fulfillments))
		return
	}

}

func TestFormatDatesForC7(t *testing.T) {
	testParams := []string{"01/02/2006 15:04", "01/02/2025 15:04", "01/02/2006", ""}
	expected := []string{"2006-01-02T15:04:00.000Z", "2025-01-02T15:04:00.000Z", "01/02/2006", ""}

	layout := "01/02/2006 15:04"

	for i, param := range testParams {
		result, err := FormatDatesForC7(layout, param)
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
		fulfillmentIds, err := GetFulfillmentIds(testCase.orderNumber, testTenant, testCase.auth, 1, nil)
		if err, ok := err.(*C7Error); ok {
			if err.StatusCode != testCase.expectedCode {
				t.Error("Test case:", i+1, "expected status code:", testCase.expectedCode, "got:", err.StatusCode)
				return
			}
		}

		if len(fulfillmentIds) > 1 {
			t.Error("fulfillmentIds length greater than 1, expected a length of 1")
			return
		}

		fulfillmentId := ""
		if len(fulfillmentIds) == 1 {
			fulfillmentId = fulfillmentIds[0]
		}

		if fulfillmentId != testCase.expectedId {
			t.Error("Test case:", i+1, "expected id:", testCase.expectedId, "got:", fulfillmentId)
			return
		}
	}
}

func Test_MarkNoFulfillmentRequired(t *testing.T) {
	orderNumber := 1020
	orderId := "b9f10447-4285-4dc2-add2-b38798dba8f9"

	// Delete previous fulfillment for test
	fulfillmentIds, err := GetFulfillmentIds(orderNumber, testTenant, AppAuthEncoded, 1, nil)
	if err != nil {
		t.Error("Error getting fulfillment id: ", err.Error())
		return
	}

	if len(fulfillmentIds) > 1 {
		t.Error("fulfillmentIds length greater than 1, expected a length of 1")
		return
	}

	fulfillmentId := fulfillmentIds[0]

	t.Log("Deleting Fulfillment ID: ", fulfillmentId)

	_, err = DeleteFulfillmentById(orderId, fulfillmentId, testTenant, AppAuthEncoded, 1, nil)
	if err != nil {
		t.Error("Error deleting fulfillment: ", err.Error())
		return
	}

	shippedTime, err := time.Parse(time.RFC3339, "2023-07-30T10:59:32.000Z")
	if err != nil {
		t.Error("Error parsing time: ", err.Error())
		return
	}

	// Mark no fulfillment required
	err = MarkNoFulfillmentRequired(orderId, shippedTime, testTenant, AppAuthEncoded, 1, nil)
	if err != nil {
		t.Error("Error marking no fulfillment required: ", err.Error())
		return
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
			return
		}
	}
}
