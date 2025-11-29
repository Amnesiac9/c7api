package c7api

import "strconv"

// Paginator is implemented by wrapper types that contain:
// - a slice of results (e.g. Orders, Products, Customers)
// - a total count
//
// Example implementers:
//   type C7OrdersFulfillmentsOnly struct {
//       Orders []C7OrderFulfillmentsOnly `json:"orders"`
//       Total  int                       `json:"total"`
//   }
type Paginator[T any] interface {
	GetItems() []T
	GetTotal() int
}

// Cursornator is implemented by wrapper types that expose:
//   - the slice of items on this "page"
//   - the cursor for the *next* page (empty string means "no more")
type Cursornator[T any] interface {
	GetItems() []T
	GetCursor() string
}

func GetAll[T any, W Paginator[T]](url string, baseQueries map[string]string, reqBody *[]byte, tenant string, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*[]T, error) {
	all := make([]T, 0, PageSize)

	// Clone the base queries so we can safely mutate page/limit
	queries := make(map[string]string, len(baseQueries)+2)
	for k, v := range baseQueries {
		queries[k] = v
	}

	page := 1
	if pStr, ok := queries["page"]; ok {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}
	for {
		queries["page"] = strconv.Itoa(page)

		wrapperPtr, err := Get[W](url, queries, reqBody, tenant, c7AppAuthEncoded, retryCount, rl)
		if err != nil {
			return nil, err
		}
		if wrapperPtr == nil {
			// treat nil as no more data
			break
		}

		wrapper := *wrapperPtr

		pageItems := wrapper.GetItems()
		total := wrapper.GetTotal()

		if len(pageItems) == 0 {
			break
		}

		all = append(all, pageItems...)

		// Stop if we already fetched all items
		if len(all) >= total {
			break
		}

		// If fewer than pageSize returned, we've hit the last page.
		if len(pageItems) < PageSize {
			break
		}

		page++
	}

	return &all, nil
}

// Cursors do not require rate limiting for now
func GetAllWithCursor[T any, W Cursornator[T]](url string, baseQueries map[string]string, reqBody *[]byte, tenant string, c7AppAuthEncoded string, retryCount int) (*[]T, error) {
	all := make([]T, 0, PageSize)

	// Clone the base queries so we can safely mutate page/limit
	queries := make(map[string]string, len(baseQueries)+2)
	for k, v := range baseQueries {
		queries[k] = v
	}

	// Start cursor unless already passed as a query
	cursor := "start"
	if c, ok := queries["cursor"]; ok {
		if c == "" {
			delete(queries, "cursor")
		} else {
			cursor = c
		}
	}
	for {
		// Only set cursor if non-empty
		if cursor != "" {
			queries["cursor"] = cursor
		}

		wrapperPtr, err := Get[W](url, queries, reqBody, tenant, c7AppAuthEncoded, retryCount, nil)
		if err != nil {
			return nil, err
		}
		if wrapperPtr == nil {
			// treat nil as no more data
			break
		}

		wrapper := *wrapperPtr

		pageItems := wrapper.GetItems()
		nextCursor := wrapper.GetCursor()

		if len(pageItems) == 0 {
			break
		}

		all = append(all, pageItems...)

		// if API returns the same cursor, avoid infinite loop
		if nextCursor == cursor {
			break
		}

		// Stop if no cursor
		if nextCursor == "" {
			break
		}

		cursor = nextCursor
	}

	return &all, nil
}
