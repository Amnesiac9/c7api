# c7api

![Tests](https://github.com/Amnesiac9/c7api/actions/workflows/Go.yml/badge.svg?branch=main)
<!-- ![Tests](https://github.com/Amnesiac9/c7api/actions/workflows/Test/badge.svg?branch=main) -->


Simple GO Package for working with the [Commerce 7](https://commerce7.com/) API.

To install:
```
github.com/Amnesiac9/c7api@latest
```

Released for personal use in my own GO projects involving Commerce7. This is a simple wrapper for making requests with built in retries and error handling.

For most requests, you can simply use the NewRequest function:
```
func NewRequest(method string, url *string, reqBody *[]byte, tenant *string, C7AppAuthEncoded *string, retryCount int) (*[]byte, error)
```
If the method is GET or DELETE, you can pass in nil for the body.

This will process a new request using the http.NewRequest, with our parameters, wrapped in a retry loop with exponential backoff. If the request returns an error or does not get a 200-299 response code, it will return a C7Error, which will include the response code and error message json from C7, if available.



