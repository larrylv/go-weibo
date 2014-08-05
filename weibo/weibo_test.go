package weibo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Weibo client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being testd.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// weibo client configured to use test server
	client = NewClient("123")
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	r.ParseForm()
	if !reflect.DeepEqual(want, r.Form) {
		t.Errorf("Request parameters = %v, want %v", r.Form, want)
	}
}

func testPostFormValues(t *testing.T, r *http.Request, values values) {
	for k, v := range values {
		if fv := r.PostFormValue(k); fv != v {
			t.Errorf("Post parameters %q is %v, want %v", k, fv, v)
		}
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("123")

	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), defaultBaseURL)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient("123")

	inURL, outURL := "foo", defaultBaseURL+"2/foo"
	req, _ := c.NewRequest("GET", inURL, nil)

	// test that relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}
}

func TestNewRequest_hasSlashPrefix(t *testing.T) {
	c := NewClient("123")

	inURL, outURL := "/foo", defaultBaseURL+"2/foo"
	req, _ := c.NewRequest("GET", inURL, nil)

	// test that relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	route := "/" + weiboApiVersion + "/"
	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	client.Do(req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	route := "/" + weiboApiVersion + "/"
	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "BadRequest", 400)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// Testing handling of an error caused by the internal http client's Do()
// function.  A redirect loop is pretty unlikely to occur within the Weibo
// API, but does allows us to exercise the right code path.
func TestDo_redirectLoop(t *testing.T) {
	setup()
	defer teardown()

	route := "/" + weiboApiVersion + "/"
	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); ok {
		t.Errorf("Expected a URL error; got %#v", err)
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusNotFound,
		Body: ioutil.NopCloser(strings.NewReader(
			`{"request": "r", "error_code": 400, "error": "e"}`,
		)),
	}
	err := CheckResponse(res)

	if err == nil {
		t.Error("Expected error response.")
	}

	want := &ErrorResponse{
		Response:   res,
		RequestURL: "r",
		ErrorCode:  400,
		Message:    "e",
	}

	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCheckResponse_noBody(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Error("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
	}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Message: "m", Response: res}
	if err.Error() == "" {
		t.Error("Expected non-empty ErrorResponse.Error()")
	}
}
