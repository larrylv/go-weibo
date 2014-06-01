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

// Status represents a Weibo's status.
type Status struct {
	ID             *int    `json:"id,omitempty"`
	Text           *string `json:"text,omitempty"`
	Source         *string `json:"source,omitempty"`
	Favorited      *bool   `json:"favorited,omitempty"`
	Truncated      *bool   `json:"truncated,omitempty"`
	RepostsCount   *int    `json:"reposts_count,omitempty"`
	CommentsCount  *int    `json:"comments_count,omitempty"`
	AttitudesCount *int    `json:"attitudes_count,omitemtpy"`
	Visible        *int    `json:"visible,omitempty"`
}

// Timeline represents Weibo statuses set.
type Timeline struct {
	Statuses       []Status `json:"statuses,omitempty"`
	TotalNumber    *int     `json:"total_number,omitempty"`
	PreviousCursor *int     `json:"previous_cursor,omitempty"`
	NextCursor     *int     `json:"next_cursor,omitempty"`
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

// UpdateOptions specifies the optional parameters to the
// StatusService.Update method.
type UpdateOptions struct {
	Status      *string  `json:status`
	Visible     *int     `json:visible,omitempty`
	ListID      *int     `json:list_id,omitempty`
	Lat         *float64 `json:lat,omitempty`
	Long        *float64 `json:long,omitempty`
	Annotations *string  `json:annotations,omitempty`
	RealIP      *string  `json:rip,omitempty`
}

// Timeline for a user. Passing the empty string will return
// timeline for the authenticated user.
//
// Weibo API docs: http://open.weibo.com/wiki/2/statuses/user_timeline
func (s *StatusesService) UserTimeline(opt *StatusListOptions) (*Timeline, *Response, error) {
	u, err := addOptions("statuses/user_timeline.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	timeline := &Timeline{}
	resp, err := s.client.Do(req, timeline)
	if err != nil {
		return nil, resp, err
	}

	return timeline, resp, err
}

func (s *StatusesService) Update(opt *UpdateOptions) (*Status, *Response, error) {
	u := "statuses/update.json"

	req, err := s.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	status := &Status{}
	resp, err := s.client.Do(req, status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, err
}
