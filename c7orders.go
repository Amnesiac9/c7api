package c7api

import "encoding/json"

func GetOrderNumberFromId(orderId string, tenant string, auth string, attempts int, rl genericRateLimiter) (int, error) {
	url := Endpoints.Order + "/" + orderId
	resp, err := RequestWithRetryAndRead("GET", url, nil, nil, tenant, auth, attempts, rl)
	if err != nil {
		return -1, err
	}

	c7Order := C7Order_OrderNumberOnly{}
	err = json.Unmarshal(*resp, &c7Order)
	if err != nil {
		return -1, err
	}

	return c7Order.OrderNumber, nil

}

func GetOrderFromId[T any](orderId string, tenant string, auth string, attempts int, rl genericRateLimiter) (*T, error) {
	url := Endpoints.Order + "/" + orderId
	resp, err := RequestWithRetryAndRead("GET", url, nil, nil, tenant, auth, attempts, rl)
	if err != nil {
		return nil, err
	}

	var result T
	err = json.Unmarshal(*resp, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
