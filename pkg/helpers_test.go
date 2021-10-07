package pkg_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/asahnoln/event-notifier/pkg"
	"google.golang.org/api/calendar/v3"
)

func paramsChecker(vals *url.Values) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		*vals = r.Form
	})
}

func fakeResponse(events []pkg.Event) http.HandlerFunc {
	items := []*calendar.Event{}
	for _, e := range events {
		items = append(items, &calendar.Event{
			Attendees: []*calendar.EventAttendee{
				{
					Email: "ivan@gmail.com",
				},
			},
			Location: e.Where,
			Summary:  e.What,
			Start: &calendar.EventDateTime{
				// TODO: Parse times for real checking
				DateTime: "2014-05-01T11:00:00Z",
			},
			End: &calendar.EventDateTime{
				// TODO: Parse times for real checking
				DateTime: "2014-05-01T13:00:00Z",
			},
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &calendar.Events{
			Items: items,
		}

		b, err := resp.MarshalJSON()
		if err != nil {
			http.Error(w, "unable to marshal response: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(b)
	})
}

func fakeServer(h http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(h)

}

func fakeTgServer(vals *url.Values) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		*vals = r.PostForm
	}))

}
