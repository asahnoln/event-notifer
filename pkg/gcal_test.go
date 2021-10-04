package pkg_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/asahnoln/event-notifier/pkg"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func TestGCalSuccessfulConnection(t *testing.T) {
	want := pkg.Event{
		"Pair Training",
		"Dom 36",
		"01.05.2014 11:00",
		"01.05.2014 13:00",
		[]string{"Andrey Kolosov"},
	}
	ts := fakeServer(fakeResponse(&want))
	defer ts.Close()

	gcal := pkg.NewGCalStore("calTestId", option.WithoutAuthentication(), option.WithEndpoint(ts.URL))
	es, err := pkg.TomorrowEvents(gcal)

	assertNoError(t, err, "unexpected error while working with gcal store: %v")
	assertSameLength(t, 1, len(es))
	assertSameStruct(t, want, es[0])
}

func TestGCalProperTiming(t *testing.T) {
	now := time.Now()
	startTomorrow := now.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	endTomorrow := startTomorrow.AddDate(0, 0, 1)
	endToday := now.AddDate(0, 0, 1).Truncate(24 * time.Hour)

	tests := []struct {
		name       string
		start, end time.Time
		eventFunc  func(store pkg.Store) ([]pkg.Event, error)
	}{
		{"tomorrow", startTomorrow, endTomorrow, pkg.TomorrowEvents},
		{"today", now, endToday, pkg.TodayEvents},
	}

	for _, tt := range tests {
		// TODO: Format time in test name so that it's easier to read
		t.Run(fmt.Sprintf("%s from %s til %s", tt.name, tt.start, tt.end), func(t *testing.T) {
			var vals url.Values

			ts := fakeServer(timeChecker(tt.start, tt.end, &vals))
			defer ts.Close()

			gcal := pkg.NewGCalStore("calTestId", option.WithoutAuthentication(), option.WithEndpoint(ts.URL))
			_, _ = tt.eventFunc(gcal)

			assertSameString(t, tt.start.Format(time.RFC3339), vals.Get("timeMin"), "want start time %q, got %q")
			assertSameString(t, tt.end.Format(time.RFC3339), vals.Get("timeMax"), "want end time %q, got %q")
		})
	}
}

func TestGCalErrorWhenWrongSettings(t *testing.T) {
	gcal := pkg.NewGCalStore("")
	_, err := pkg.TomorrowEvents(gcal)
	assertError(t, err, "expected error with no calendar id, got nil")

}

func timeChecker(start, end time.Time, vals *url.Values) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		*vals = r.Form
	})
}

func fakeResponse(data *pkg.Event) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &calendar.Events{
			Items: []*calendar.Event{
				{
					Attendees: []*calendar.EventAttendee{
						{
							DisplayName: data.Who[0],
						},
					},
					Location: data.Where,
					Summary:  data.What,
					Start: &calendar.EventDateTime{
						DateTime: "2014-05-01T11:00:00Z",
					},
					End: &calendar.EventDateTime{
						DateTime: "2014-05-01T13:00:00Z",
					},
				},
			},
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

func assertNoError(t testing.TB, err error, message string) {
	t.Helper()

	if err != nil {
		t.Fatalf(message, err)
	}
}

func assertError(t testing.TB, err error, message string) {
	t.Helper()

	if err == nil {
		t.Fatal(message)
	}
}
