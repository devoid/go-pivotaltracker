// Copyright (C) 2015 Scott Devoid
// Use of this source code is governed by the MIT License.
// The license can be found in the LICENSE file.

package pivotal

import (
	"net/http"
)

// StoryServiceShim has the same functions as StoryService but
// preserves calling convention where the project ID is included
// in each function call. Also preserves old filter by string on
// the List() function.
type StoryServiceShim struct {
	*Client
}

func newStoryServiceShim(c *Client) *StoryServiceShim {
	return &StoryServiceShim{c}
}

func (s *StoryServiceShim) List(projectId int, filter string) ([]*Story, *http.Response, error) {
	return s.Client.Project(projectId).Stories.List(Filter(filter))
}

func (s *StoryServiceShim) Get(projectId, storyId int) (*Story, *http.Response, error) {
	return s.Client.Project(projectId).Stories.Get(storyId)
}

func (s *StoryServiceShim) Update(projectId, storyId int, story *Story) (*Story, *http.Response, error) {
	return s.Client.Project(projectId).Stories.Update(storyId, story)
}

func (s *StoryServiceShim) ListTasks(projectId, storyId int) ([]*Task, *http.Response, error) {
	return s.Client.Project(projectId).Stories.ListTasks(storyId)
}

func (s *StoryServiceShim) AddTask(projectId, storyId int, task *Task) (*http.Response, error) {
	return s.Client.Project(projectId).Stories.AddTask(storyId, task)
}

func (s *StoryServiceShim) ListOwners(projectId, storyId int) ([]*Person, *http.Response, error) {
	return s.Client.Project(projectId).Stories.ListOwners(storyId)
}

func (s *StoryServiceShim) AddComment(projectId, storyId int, comment *Comment) (*Comment, *http.Response, error) {
	return s.Client.Project(projectId).Stories.AddComment(storyId, comment)
}
