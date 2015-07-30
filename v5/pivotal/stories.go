// Copyright (c) 2014 Salsita Software
// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.
package pivotal

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	StoryTypeFeature = "feature"
	StoryTypeBug     = "bug"
	StoryTypeChore   = "chore"
	StoryTypeRelease = "release"
)

const (
	StoryStateUnscheduled = "unscheduled"
	StoryStatePlanned     = "planned"
	StoryStateUnstarted   = "unstarted"
	StoryStateStarted     = "started"
	StoryStateFinished    = "finished"
	StoryStateDelivered   = "delivered"
	StoryStateAccepted    = "accepted"
	StoryStateRejected    = "rejected"
)

type Story struct {
	Id            int        `json:"id,omitempty"`
	ProjectId     int        `json:"project_id,omitempty"`
	Name          string     `json:"name,omitempty"`
	Description   string     `json:"description,omitempty"`
	Type          string     `json:"story_type,omitempty"`
	State         string     `json:"current_state,omitempty"`
	Estimate      *float64   `json:"estimate,omitempty"`
	AcceptedAt    *time.Time `json:"accepted_at,omitempty"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	RequestedById int        `json:"requested_by_id,omitempty"`
	OwnerIds      *[]int     `json:"owner_ids,omitempty"`
	LabelIds      *[]int     `json:"label_ids,omitempty"`
	Labels        *[]*Label  `json:"labels,omitempty"`
	TaskIds       *[]int     `json:"task_ids,omitempty"`
	Tasks         *[]int     `json:"tasks,omitempty"`
	FollowerIds   *[]int     `json:"follower_ids,omitempty"`
	CommentIds    *[]int     `json:"comment_ids,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	IntegrationId int        `json:"integration_id,omitempty"`
	ExternalId    string     `json:"external_id,omitempty"`
	URL           string     `json:"url,omitempty"`
	Kind          string     `json:"kind,omitempty"`
}

type Task struct {
	Id          int        `json:"id,omitempty"`
	StoryId     int        `json:"story_id,omitempty"`
	Description string     `json:"description,omitempty"`
	Position    int        `json:"position,omitempty"`
	Complete    bool       `json:"complete,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type Person struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Initials string `json:"initials,omitempty"`
	Username string `json:"username,omitempty"`
	Kind     string `json:"kind,omitempty"`
}

type Comment struct {
	Id                  int        `json:"id,omitempty"`
	StoryId             int        `json:"story_id,omitempty"`
	EpicId              int        `json:"epic_id,omitempty"`
	PersonId            int        `json:"person_id,omitempty"`
	Text                string     `json:"text,omitempty"`
	FileAttachmentIds   []int      `json:"file_attachment_ids,omitempty"`
	GoogleAttachmentIds []int      `json:"google_attachment_ids,omitempty"`
	CommitType          string     `json:"commit_type,omitempty"`
	CommitIdentifier    string     `json:"commit_identifier,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
}

type StoryService struct {
	client    *Client
	projectId string
}

func newStoryService(client *Client, projectId string) *StoryService {
	return &StoryService{client, projectId}
}

func (s *StoryService) setupReq(opts []RequestOption) (req *http.Request, err error) {
	u := fmt.Sprintf("projects/%v/stories", s.projectId)
	req, err = s.client.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	for _, opt := range opts {
		opt(req)
	}
	return
}

func (s *StoryService) List(opts ...RequestOption) ([]*Story, *http.Response, error) {
	req, err := s.setupReq(opts)
	if err != nil {
		return nil, nil, err
	}

	var stories []*Story
	resp, err := s.client.Do(req, &stories)
	if err != nil {
		return nil, resp, err
	}

	return stories, resp, err
}

type StoryCursor struct {
	*cursor
	buff []*Story
}

func (c *StoryCursor) Next() (s *Story, err error) {
	if len(c.buff) == 0 {
		_, err = c.next(&c.buff)
		if err != nil {
			return nil, err
		}
	}

	if len(c.buff) == 0 {
		err = io.EOF
	} else {
		s, c.buff = c.buff[0], c.buff[1:]
	}
	return s, err
}

func (s *StoryService) Iterate(opts ...RequestOption) (c *StoryCursor, err error) {
	req_fn := func() (req *http.Request) {
		req, _ = s.setupReq(opts)
		return req
	}
	cc, err := newCursor(s.client, req_fn)
	return &StoryCursor{cc, make([]*Story, 0)}, err
}

func (service *StoryService) Get(storyId int) (*Story, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/stories/%v", service.projectId, storyId)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var story Story
	resp, err := service.client.Do(req, &story)
	if err != nil {
		return nil, resp, err
	}

	return &story, resp, err
}

func (service *StoryService) Update(storyId int, story *Story) (*Story, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/stories/%v", service.projectId, storyId)
	req, err := service.client.NewRequest("PUT", u, story)
	if err != nil {
		return nil, nil, err
	}

	var bodyStory Story
	resp, err := service.client.Do(req, &bodyStory)
	if err != nil {
		return nil, resp, err
	}

	return &bodyStory, resp, err

}

func (service *StoryService) ListTasks(storyId int) ([]*Task, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/stories/%v/tasks", service.projectId, storyId)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tasks []*Task
	resp, err := service.client.Do(req, &tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (service *StoryService) AddTask(storyId int, task *Task) (*http.Response, error) {
	if task.Description == "" {
		return nil, &ErrFieldNotSet{"description"}
	}

	u := fmt.Sprintf("projects/%v/stories/%v/tasks", service.projectId, storyId)
	req, err := service.client.NewRequest("POST", u, task)
	if err != nil {
		return nil, err
	}

	return service.client.Do(req, nil)
}

func (service *StoryService) ListOwners(storyId int) ([]*Person, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/stories/%d/owners", service.projectId, storyId)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var owners []*Person
	resp, err := service.client.Do(req, &owners)
	if err != nil {
		return nil, resp, err
	}

	return owners, resp, err
}

func (service *StoryService) AddComment(storyId int, comment *Comment) (*Comment, *http.Response, error) {
	u := fmt.Sprintf("projects/%v/stories/%v/comments", service.projectId, storyId)
	req, err := service.client.NewRequest("POST", u, comment)
	if err != nil {
		return nil, nil, err
	}

	var newComment Comment
	resp, err := service.client.Do(req, &newComment)
	if err != nil {
		return nil, resp, err
	}

	return &newComment, resp, err
}
