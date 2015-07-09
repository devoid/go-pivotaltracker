// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"fmt"
	"io"
	"net/http"
	//	"strconv"
	"time"
)

type Iteration struct {
	Number       int        `json:"number,omitempty"`
	ProjectId    int        `json:"project_id,omitempty"`
	Length       int        `json:"length,omitempty"`
	TeamStrength float64    `json:"team_strength,omitempty"`
	StoryIds     []int      `json:"story_ids,omitempty"`
	Start        *time.Time `json:"start,omitempty"`
	Finish       *time.Time `json:"finish,omitempty"`
	Kind         string     `json:"kind,omitempty"`
}

type IterationOverride struct {
	Number       int     `json:"number,omitempty"`
	ProjectId    int     `json:"project_id,omitempty"`
	Length       int     `json:"length,omitempty"`
	TeamStrength float64 `json:"team_strength,omitempty"`
	Kind         string  `json:"kind,omitempty"`
}

type IterationOverrideRequest struct {
	IterationNumber int     `json:"iteration_number,omitempty"`
	ProjectId       int     `json:"project_id,omitempty"`
	Length          int     `json:"length,omitempty"`
	TeamStrength    float64 `json:"team_strength,omitempty"`
}

type IterationService struct {
	*Client
	projectId string
}

func newIterationService(client *Client, projectId string) *IterationService {
	return &IterationService{client, projectId}
}

func (s *IterationService) List(opts ...RequestOption) ([]*Iteration, *http.Response, error) {
	req_fn := func() (req *http.Request) {
		req, _ = s.setupReq(opts...)
		return req
	}
	cc, err := newCursor((s).Client, req_fn, 10)
	if err != nil {
		return nil, nil, err
	}
	var iterations []*Iteration
	var resp *http.Response
	for {
		var its []*Iteration
		resp, err = cc.next(&its)
		if err != nil && err != io.EOF {
			break
		}
		iterations = append(iterations, its...)
		if err == io.EOF {
			break
		}
	}
	if err == io.EOF {
		err = nil
	}
	return iterations, resp, err
}

type IterationCursor struct {
	*cursor
	buff []*Iteration
}

func (c *IterationCursor) Next() (i *Iteration, err error) {
	if len(c.buff) == 0 {
		_, err = c.next(&c.buff)
		if err != nil {
			return
		}
	}
	if len(c.buff) == 0 {
		err = io.EOF
	} else {
		i, c.buff = c.buff[0], c.buff[1:]
	}
	return
}

func (s *IterationService) Iterate(opts ...RequestOption) (c *IterationCursor, err error) {
	req_fn := func() (req *http.Request) {
		req, _ = s.setupReq(opts...)
		return req
	}
	cc, err := newCursor((s).Client, req_fn, 10)
	return &IterationCursor{cc, make([]*Iteration, 0)}, err
}

func (s *IterationService) OverrideIteration(o IterationOverrideRequest) (
	*IterationOverride, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/iterations/%v", s.projectId, o.IterationNumber)
	req, err := s.NewRequest("PUT", u, o)
	if err != nil {
		return nil, nil, err
	}
	var iteration *IterationOverride
	resp, err := s.Do(req, &iteration)
	if err != nil {
		return nil, resp, err
	}
	return iteration, resp, err
}

func (s *IterationService) setupReq(opts ...RequestOption) (req *http.Request, err error) {
	u := fmt.Sprintf("projects/%v/iterations", s.projectId)
	req, err = s.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	for _, opt := range opts {
		opt(req)
	}
	return
}
