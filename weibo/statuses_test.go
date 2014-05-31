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
		fmt.Fprint(w, `[{"id": "1"}]`)
	})

	opt := &StatusListOptions{UID: String(uid)}
	statuses, _, err := client.Statuses.UserTimeline(opt)

	if err != nil {
		t.Errorf("Statuses.UserTimeline returned error: %v", err)
	}

	want := []Status{{ID: String("1")}}
	if !reflect.DeepEqual(statuses, want) {
		t.Errorf("Statuses.UserTimeline returned %+v, want %+v", statuses, want)
	}
}
