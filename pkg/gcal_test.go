package pkg_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
	ts := fakeServer(&want)
	defer ts.Close()

	gcal := pkg.NewGCalStore("calTestId", option.WithoutAuthentication(), option.WithEndpoint(ts.URL))
	es, err := pkg.TomorrowEvents(gcal)

	assertNoError(t, err, "unexpected error while working with gcal store: %v")
	assertSameLength(t, 1, len(es))
	assertSameStruct(t, want, es[0])
}

func TestGCalErrorWhenWrongSettings(t *testing.T) {
	gcal := pkg.NewGCalStore("")
	_, err := pkg.TomorrowEvents(gcal)
	assertError(t, err, "expected error with no calendar id, got nil")

}

func fakeServer(data *pkg.Event) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))

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
