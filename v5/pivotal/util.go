// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"io"
	"net/http"
	"strconv"
	"sync"
)

// requestFn is a function that returns a new *http.Request object.
type requestFn func() (req *http.Request)

// cursor tracks response headers from paginated API responses.
// And sets the appropriate URI variables in the next request.
type cursor struct {
	client    *Client
	requestFn requestFn
	limit     int
	offset    int
	reqCount  int
	lock      *sync.Mutex
}

func newCursor(client *Client, fn requestFn) (c *cursor, err error) {
	// Default to 10 items, which seems to be what Pivotal natively returns.
	return &cursor{
		client:    client,
		requestFn: fn,
		limit:     10,
		lock:      &sync.Mutex{},
	}, nil
}

// next is called with a pointer to an []*Type, which will be correctly
// unmarshalled. next() returns the http.Response, where the Body is
// already closed and an error. When next() reaches the end of a paginated
// endpoint, it returns io.EOF as the error.
func (c *cursor) next(v interface{}) (resp *http.Response, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	req := c.requestFn()

	// Note: if we've already made requests, we always update
	// the limit and offset fields to the current value. If we've
	// never made requests the user may have set these values
	// in requestFn; So we honor that decision the first time.
	values := req.URL.Query()
	if c.reqCount != 0 {
		values.Set("limit", strconv.Itoa(c.limit))
		values.Set("offset", strconv.Itoa(c.offset))

	} else {
		if values.Get("limit") == "" {
			values.Set("limit", strconv.Itoa(c.limit))
		}
		if values.Get("offset") == "" {
			values.Set("offset", strconv.Itoa(c.offset))
		}
	}
	req.URL.RawQuery = values.Encode()

	// Do the request, decode JSON to v
	resp, err = c.client.Do(req, &v)
	// increment request counter
	c.reqCount++
	if err != nil {
		return nil, err
	}

	// Helper to extract and convert Header values that are Int's
	getIntHeader := func(resp *http.Response, k string) int {
		if err != nil {
			return 0
		}
		i, cerr := strconv.Atoi(resp.Header.Get(k))
		if cerr != nil {
			err = cerr
		}
		return i
	}

	// Get limit, offset, total and returned headers for pagination
	limit := getIntHeader(resp, "X-Tracker-Pagination-Limit")
	offset := getIntHeader(resp, "X-Tracker-Pagination-Offset")
	total := getIntHeader(resp, "X-Tracker-Pagination-Total")
	returned := getIntHeader(resp, "X-Tracker-Pagination-Returned")
	if err != nil {
		return nil, err
	}

	// Calculate the new offset, which is the old offset plus
	// the minimum of (returned, limit)
	if returned < limit {
		c.offset = offset + returned
	} else {
		c.offset = offset + limit
	}

	// Return EOF if we have reached the end.
	if c.offset >= total {
		err = io.EOF
	}
	return resp, err
}
