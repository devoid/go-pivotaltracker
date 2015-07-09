// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"fmt"
	"net/http"
	"time"
)

type Label struct {
	Id        int        `json:"id,omitempty"`
	ProjectId int        `json:"project_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Kind      string     `json:"kind,omitempty"`
}

type LabelService struct {
	client    *Client
	projectId string
}

func newLabelService(client *Client, projectId string) *LabelService {
	return &LabelService{client, projectId}
}

func (s *LabelService) List() (labels []*Label, resp *http.Response, err error) {
	u := fmt.Sprintf("projects/%s/labels", s.projectId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, &labels)
	if err != nil {
		return nil, resp, err
	}
	return
}

func (s *LabelService) Create(name string) (label *Label, resp *http.Response, err error) {
	l := Label{Name: name}
	u := fmt.Sprintf("projects/%s/labels", s.projectId)
	req, err := s.client.NewRequest("POST", u, l)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, &label)
	if err != nil {
		return nil, resp, err
	}
	return
}

func (s *LabelService) Get(id int) (label *Label, resp *http.Response, err error) {
	u := fmt.Sprintf("projects/%s/labels/%d", s.projectId, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, &label)
	if err != nil {
		return nil, resp, err
	}
	return
}

func (s *LabelService) Rename(id int, newName string) (
	label *Label, resp *http.Response, err error) {
	l := Label{Name: newName}
	u := fmt.Sprintf("projects/%s/labels/%d", s.projectId, id)
	req, err := s.client.NewRequest("PUT", u, l)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, &label)
	if err != nil {
		return nil, resp, err
	}
	return
}

func (s *LabelService) Delete(id int) (resp *http.Response, err error) {
	u := fmt.Sprintf("projects/%s/labels/%d", s.projectId, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	return
}
