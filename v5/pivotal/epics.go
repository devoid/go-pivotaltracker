// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Epic struct {
	Id          int        `json:"id,omitempty"`
	ProjectId   int        `json:"project_id,omitempty"`
	Name        string     `json:"name,omitempty"`
	LabelId     int        `json:"label_id,omitempty"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Kind        string     `json:"kind,omitempty"`
}

type EpicRequest struct {
	ProjectId   int        `json:"project_id,omitempty"`
	Name        string     `json:"name,omitempty"`
	LabelId     int        `json:"label_id,omitempty"`
	Description string     `json:"description,omitempty"`
	Comments    []*Comment `json:"comments,omitempty"`
	Followers   []*Person  `json:"followers,omitempty"`
	FollowerIds []int      `json:"follower_ids,omitempty"`
	AfterId     int        `json:"after_id,omitempty"`
	BeforeId    int        `json:"before_id,omitempty"`
}

type EpicService struct {
	*Client
	projectId string
}

func newEpicService(client *Client, projectId string) *EpicService {
	return &EpicService{client, projectId}
}

func (s *EpicService) List(opts ...RequestOption) ([]*Epic, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/epics", s.projectId)
	req, err := s.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	for _, opt := range opts {
		opt(req)
	}
	var epics []*Epic
	resp, err := s.Do(req, &epics)
	if err != nil {
		return nil, resp, err
	}
	return epics, resp, err
}

func (s *EpicService) Get(id int) (*Epic, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/epics/%v", s.projectId, id)
	req, err := s.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	var epic *Epic
	resp, err := s.Do(req, &epic)
	if err != nil {
		return nil, resp, err
	}
	return epic, resp, err
}

func (s *EpicService) Add(epic EpicRequest) (*Epic, *http.Response, error) {
	project := s.projectId
	if epic.ProjectId != 0 && strconv.Itoa(epic.ProjectId) != project {
		project = strconv.Itoa(epic.ProjectId)
	}
	u := fmt.Sprintf("projects/%v/epics", project)
	req, err := s.NewRequest("POST", u, epic)
	if err != nil {
		return nil, nil, err
	}
	var e *Epic
	resp, err := s.Do(req, &e)
	if err != nil {
		return nil, resp, err
	}
	return e, resp, err
}

func (s *EpicService) Update(epicId int, epic EpicRequest) (*Epic, *http.Response, error) {
	project := s.projectId
	if epic.ProjectId != 0 && strconv.Itoa(epic.ProjectId) != project {
		project = strconv.Itoa(epic.ProjectId)
	}
	u := fmt.Sprintf("projects/%v/epics/%v", project, epicId)
	req, err := s.NewRequest("PUT", u, epic)
	if err != nil {
		return nil, nil, err
	}
	var e *Epic
	resp, err := s.Do(req, &e)
	if err != nil {
		return nil, resp, err
	}
	return e, resp, err
}

func (s *EpicService) Delete(epicId int) (resp *http.Response, err error) {
	u := fmt.Sprintf("projects/%v/epics/%v", s.projectId, epicId)
	req, err := s.NewRequest("DELETE", u, nil)
	if err != nil {
		return
	}
	return s.Do(req, nil)
}
