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
	CreatedAt      *string  `json:"created_at,omitempty"`
	ID             *int64   `json:"id,omitempty"`
	MID            *string  `json:"mid,omitempty"`
	IDStr          *string  `json:"idstr,omitempty"`
	Text           *string  `json:"text,omitempty"`
	Source         *string  `json:"source,omitempty"`
	Favorited      *bool    `json:"favorited,omitempty"`
	Truncated      *bool    `json:"truncated,omitempty"`
	User           *User    `json:"user,omitempty"`
	RepostsCount   *int     `json:"reposts_count,omitempty"`
	CommentsCount  *int     `json:"comments_count,omitempty"`
	AttitudesCount *int     `json:"attitudes_count,omitemtpy"`
	Visible        *Visible `json:"visible,omitempty"`
}

// Visible represents visible object of a Weibo status.
type Visible struct {
	VType  *int `json:"type,omitempty"`
	ListID *int `json:"list_id,omitempty"`
}

// Timeline represents Weibo statuses set.
type Timeline struct {
	Statuses       []Status `json:"statuses,omitempty"`
	TotalNumber    *int     `json:"total_number,omitempty"`
	PreviousCursor *int     `json:"previous_cursor,omitempty"`
	NextCursor     *int     `json:"next_cursor,omitempty"`
}

// TimelineIDs represents Weibo statuses ids set.
type TimelineIDs struct {
	StatusesIDs    []string `json:"statuses,omitempty"`
	TotalNumber    *int     `json:"total_number,omitempty"`
	PreviousCursor *int     `json:"previous_cursor,omitempty"`
	NextCursor     *int     `json:"next_cursor,omitempty"`
}

// StatusListOptions specifies the optional parameters to the
// StatusService.UserTimeline method.
type StatusListOptions struct {
	UID        string `url:"uid,omitempty"`
	ScreenName string `url:"screen_name,omitempty"`
	SinceID    string `url:"since_id,omitempty"`
	MaxID      string `url:"max_id,omitempty"`
	ListOptions
}

// StatusRequest represetns a request to create a status.
type StatusRequest struct {
	Status      *string  `url:"status"`
	Visible     *int     `url:"visible,omitempty"`
	ListID      *int     `url:"list_id,omitempty"`
	Lat         *float64 `url:"lat,omitempty"`
	Long        *float64 `url:"long,omitempty"`
	Annotations *string  `url:"annotations,omitempty"`
	RealIP      *string  `url:"rip,omitempty"`
}

// Timeline of a user. Passing the empty string will return
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

// Timeline IDs of a user. Passing the empty string will return
// timeline for the authenticated user.
//
// Weibo API docs: http://open.weibo.com/wiki/2/statuses/user_timeline
func (s *StatusesService) UserTimelineIDs(opt *StatusListOptions) (*TimelineIDs, *Response, error) {
	u, err := addOptions("statuses/user_timeline/ids.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	timelineIDs := &TimelineIDs{}
	resp, err := s.client.Do(req, timelineIDs)
	if err != nil {
		return nil, resp, err
	}

	return timelineIDs, resp, err
}

// Create a Weibo Status.
//
// Weibo API docs: http://open.weibo.com/wiki/2/statuses/update
func (s *StatusesService) Create(opt *StatusRequest) (*Status, *Response, error) {
	u := "statuses/update.json"

	req, err := s.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	status := new(Status)
	resp, err := s.client.Do(req, status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, err
}
