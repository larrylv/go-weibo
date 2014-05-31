package weibo

import (
	_ "fmt"
)

// StatusesService handles communication with the Status related
// methods of the Weibo API.
//
// Weibo API docs: http://open.weibo.com/wiki/%E5%BE%AE%E5%8D%9AAPI
type StatusesService struct {
	client *Client
}

// Status represents a Weibo's status
type Status struct {
	ID             *string `json:"id,omitempty"`
	Text           *string `json:"text,omitempty"`
	Source         *string `json:"source,omitempty"`
	Favorited      *bool   `json:"favorited,omitempty"`
	Truncated      *bool   `json:"truncated,omitempty"`
	RepostsCount   *int    `json:"reposts_count,omitempty"`
	CommentsCount  *int    `json:"comments_count,omitempty"`
	AttitudesCount *int    `json:"attitudes_count,omitemtpy"`
	Visible        *int    `json:"visible,omitempty"`
}

// StatusListOptions specifies the optional parameters to the
// StatusService.UserTimeline method.
type StatusListOptions struct {
	UID        *string `url:"uid,omitempty"`
	ScreenName *string `url:"screen_name,omitempty"`
	SinceID    *string `url:"since_id,omitempty"`
	MaxID      *string `url:"max_id,omitempty"`
	ListOptions
}

// Timeline for a user. Passing the empty string will return
// timeline for the authenticated user.
//
// Weibo API docs: http://open.weibo.com/wiki/2/statuses/user_timeline
func (s *StatusesService) UserTimeline(opt *StatusListOptions) ([]Status, *Response, error) {
	u, err := addOptions("statuses/user_timeline", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	statuses := new([]Status)
	resp, err := s.client.Do(req, statuses)
	if err != nil {
		return nil, resp, err
	}

	return *statuses, resp, err
}
