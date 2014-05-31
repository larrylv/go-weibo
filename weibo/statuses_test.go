package weibo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUserTimeline(t *testing.T) {
	setup()
	defer teardown()

	uid := "42"

	mux.HandleFunc("/2/statuses/user_timeline", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"uid": uid,
		})
		fmt.Fprint(w, `{"statuses": [{"id": 1, "text": "hello weibo"}], "total_number": 1}`)
	})

	opt := &StatusListOptions{UID: String(uid)}
	timeline, _, err := client.Statuses.UserTimeline(opt)

	if err != nil {
		t.Errorf("Statuses.UserTimeline returned error: %v", err)
	}

	want := Timeline{Statuses: []Status{{ID: Int(1), Text: String("hello weibo")}}, TotalNumber: Int(1)}
	if !reflect.DeepEqual(timeline.Statuses, want.Statuses) {
		t.Errorf("Statuses.UserTimeline returned %+v, want %+v", timeline, want)
	}
}
