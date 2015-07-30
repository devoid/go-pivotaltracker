// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"net/http"
	"strconv"
	"time"
)

type RequestOption func(*http.Request)

func Limit(i int) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "limit", strconv.Itoa(i))
	}
}

func Offset(i int) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "offset", strconv.Itoa(i))
	}
}

func WithLabel(label string) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "with_label", label)
	}
}

func WithState(state string) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "with_state", state)
	}
}

func AfterStory(id int) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "after_story_id", strconv.Itoa(id))
	}
}

func BeforeStory(id int) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "before_story_id", strconv.Itoa(id))
	}
}

func AcceptedBefore(t *time.Time) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "accepted_before", t.String())
	}
}

func CreatedBefore(t *time.Time) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "created_before", t.String())
	}
}

func CreatedAfter(t *time.Time) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "created_after", t.String())
	}
}

func Filter(s string) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "query", s)
	}
}

func WithScope(scope string) RequestOption {
	return func(r *http.Request) {
		urlAddParam(r, "scope", scope)
	}
}

func urlAddParam(r *http.Request, k, v string) {
	query := r.URL.Query()
	query.Add(k, v)
	r.URL.RawQuery = query.Encode()
}
