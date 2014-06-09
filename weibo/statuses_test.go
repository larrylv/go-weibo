package weibo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestStatusesUserTimeline(t *testing.T) {
	setup()
	defer teardown()

	uid := "42"

	mux.HandleFunc("/2/statuses/user_timeline.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"uid": uid,
		})
		fmt.Fprint(w, `{"statuses": [{"id": 1, "text": "hello weibo", "user": {"id": 42, "name": "larrylv"}}], "total_number": 1}`)
	})

	opt := &StatusListOptions{UID: uid}
	timeline, _, err := client.Statuses.UserTimeline(opt)

	if err != nil {
		t.Errorf("Statuses.UserTimeline returned error: %v", err)
	}

	want := Timeline{Statuses: []Status{{ID: Int64(1), Text: String("hello weibo"), User: &User{ID: Int(42), Name: String("larrylv")}}}, TotalNumber: Int(1)}
	if !reflect.DeepEqual(timeline.Statuses, want.Statuses) {
		t.Errorf("Statuses.UserTimeline returned %+v, want %+v", timeline, want)
	}
}

func TestStatusesUserTimelineIDs(t *testing.T) {
	setup()
	defer teardown()

	uid := "42"

	mux.HandleFunc("/2/statuses/user_timeline/ids.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"uid": uid,
		})
		fmt.Fprint(w, `{"statuses": ["1234", "5678"], "total_number": 2}`)
	})

	opt := &StatusListOptions{UID: uid}
	timelineIDs, _, err := client.Statuses.UserTimelineIDs(opt)

	if err != nil {
		t.Errorf("Statuses.UserTimeline returned error: %v", err)
	}

	want := TimelineIDs{StatusesIDs: []string{"1234", "5678"}, TotalNumber: Int(2)}
	if !reflect.DeepEqual(timelineIDs.StatusesIDs, want.StatusesIDs) {
		t.Errorf("Statuses.UserTimelineIDs returned %+v, want %+v", timelineIDs, want)
	}
}

func TestStatusesUpdate(t *testing.T) {
	setup()
	defer teardown()

	text := "Hello, weibo!"

	opt := &StatusRequest{
		Status:  String(text),
		Visible: Int(1),
	}

	mux.HandleFunc("/2/statuses/update.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testPostFormValues(t, r, values{
			"status":  text,
			"visible": "1",
		})

		fmt.Fprint(w,
			`
            {
                "id" : 1,
                "text" : "Hello, weibo!",
                "visible": {
                    "type": 1
                }
            }
            `)
	})

	status, _, err := client.Statuses.Create(opt)

	if err != nil {
		t.Errorf("Statuses.Update returned error %v", err)
	}

	want := &Status{ID: Int64(1), Text: String("Hello, weibo!"), Visible: &Visible{VType: Int(1)}}
	if !reflect.DeepEqual(status, want) {
		t.Errorf("Statuses.Update returned %+v, want %+v", status, want)
	}
}
