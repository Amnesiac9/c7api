# c7api

![Tests](https://github.com/Amnesiac9/c7api/actions/workflows/tests.yml/badge.svg?branch=main)
![Build](https://github.com/Amnesiac9/c7api/actions/workflows/build.yml/badge.svg?branch=main)


Simple GO Package for working with the [Commerce 7](https://commerce7.com/) API.

CURRENTLY IN DEVELOPMENT - EXPECT BREAKING CHANGES

To install:
```
github.com/Amnesiac9/c7api@latest
```

Released for personal use in my own GO projects involving Commerce7. This is a simple wrapper for making requests with built in retries and error handling.

For most requests, you can simply use the NewRequest function:
```
Request(method string, url string, reqBody *[]byte, tenant string, c7AppAuthEncoded string) (*http.Response, error)
```
If the method is GET or DELETE, you can pass in nil for the body.

For requests with backoff: 
```
RequestWithRetryAndRead(method string, url string, reqBody *[]byte, tenant string, c7AppAuthEncoded string, retryCount int) (*[]byte, error)
```
This will process a new request using the http.NewRequest, with our parameters, wrapped in a retry loop with exponential backoff. If the request returns an error or does not get a 200-299 response code, it will return a C7Error, which will include the response code and error message json from C7, if available.
