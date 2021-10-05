package pkg_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/event-notifier/pkg"
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

	gcal := pkg.NewGCalStore(
		"calTestId",
		strings.NewReader(`{"ivan@gmail.com": "Andrey Kolosov"}`),
		option.WithoutAuthentication(),
		option.WithEndpoint(ts.URL),
	)
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

			ts := fakeServer(paramsChecker(&vals))
			defer ts.Close()

			gcal := pkg.NewGCalStore("calTestId", nil, option.WithoutAuthentication(), option.WithEndpoint(ts.URL))
			_, _ = tt.eventFunc(gcal)

			assertSameString(t, tt.start.Format(time.RFC3339), vals.Get("timeMin"), "want start time %q, got %q")
			assertSameString(t, tt.end.Format(time.RFC3339), vals.Get("timeMax"), "want end time %q, got %q")
		})
	}
}

func TestGCalErrorWhenWrongSettings(t *testing.T) {
	gcal := pkg.NewGCalStore("", nil)
	_, err := pkg.TomorrowEvents(gcal)
	assertError(t, err, "expected error with no calendar id, got nil")
}

func TestGCalUseProperParams(t *testing.T) {
	var vals url.Values
	ts := fakeServer(paramsChecker(&vals))
	gcal := pkg.NewGCalStore("testId", nil, option.WithoutAuthentication(), option.WithEndpoint(ts.URL))
	_, _ = pkg.TomorrowEvents(gcal)

	assertSameString(t, "false", vals.Get("showDeleted"), "want showDeleted %q, got %q")
	assertSameString(t, "true", vals.Get("singleEvents"), "want singleEvents %q, got %q")
	assertSameString(t, "startTime", vals.Get("orderBy"), "want orderBy %q, got %q")
}
