package weibo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion  = "0.1"
	weiboApiVersion = "2"
	defaultBaseURL  = "https://api.weibo.com/"
	userAgent       = "go-weibo/" + libraryVersion
)

// A Client manages communication with the Weibo API.
type Client struct {
	// HTTP client used to communcate with the API.
	client *http.Client

	// Access Token
	accessToken string

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent used when communicating with the Weibo API.
	UserAgent string

	// Services used for talking to different parts of the Weibo API.
	Statuses *StatusesService
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"count,omitempty"`
}

// addOptions adds the parameters in opt as URL query parameters to string.
// opt must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewClient returns a new Weibo API client.
func NewClient(accessToken string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: http.DefaultClient, accessToken: accessToken, BaseURL: baseURL, UserAgent: userAgent}
	c.Statuses = &StatusesService{client: c}

	return c
}

// NewRequest creates an API request.  A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlString string, body interface{}) (*http.Request, error) {
	if strings.HasPrefix(urlString, "/") {
		urlString = weiboApiVersion + urlString
	} else {
		urlString = weiboApiVersion + "/" + urlString
	}

	rel, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Authorization", "OAuth2 "+c.accessToken)
	return req, nil
}

// Response is a Weibo API response.
// This wraps the standrad http.Response returned from Weibo.
// For future use, maybe.
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}

	return response
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occured.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting
// to first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	if err := CheckResponse(resp); err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

// An Error Response reports one or more errors caused by an API request.
//
// Weibo API docs: http://open.weibo.com/wiki/Error_code
type ErrorResponse struct {
	Response   *http.Response // HTTP response that caused this error
	RequestURL string         `json:"request"`    // request on which the error occured
	ErrorCode  int            `json:"error_code"` // error_code
	Message    string         `json:"error"`      // error message
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.RequestURL,
		r.Response.StatusCode, r.ErrorCode, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it, but unlike Int32
// its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}
