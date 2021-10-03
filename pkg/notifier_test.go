package pkg_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/asahnoln/event-notifier/pkg"
)

type stubStore struct {
	events map[pkg.EventType][]pkg.Event
}

func (s *stubStore) Events(when pkg.EventType) ([]pkg.Event, error) {
	return s.events[when], nil
}

var db = &stubStore{
	map[pkg.EventType][]pkg.Event{
		pkg.Tomorrow: {
			{
				What:  "Scene 2.5",
				Where: "Dom, Great Hall",
				Who:   []string{"Ivan", "Erkanat"},
				Start: "08.11.2019 19:00",
				End:   "08.11.2019 22:00",
			},
		},
		pkg.Today: {
			{
				What:  "Scene 2.6",
				Where: "Dom, Second Room",
				Who:   []string{"Varvara", "Kamila"},
				Start: "04.10.2021 15:00",
				End:   "04.10.2019 18:00",
			},
		},
	},
}

func TestGetTomorrowEvents(t *testing.T) {
	es, _ := pkg.TomorrowEvents(db)
	assertSameLength(t, 1, len(es))

	assertSameStruct(t, db.events[pkg.Tomorrow][0], es[0])
}

func TestGetTodayEvents(t *testing.T) {
	es, _ := pkg.TodayEvents(db)

	assertSameStruct(t, db.events[pkg.Today][0], es[0])
}

type stubSender struct {
	result string
}

func (s *stubSender) Send(message string) error {
	s.result = message
	return nil
}

func TestSendMessageForEvent(t *testing.T) {
	es, _ := pkg.TomorrowEvents(db)

	sdr := &stubSender{}
	pkg.Send(es, sdr)

	want := es[0]
	assertContains(t, want.What, sdr.result)
	assertContains(t, want.Where, sdr.result)
	for _, p := range want.Who {
		assertContains(t, p, sdr.result)
	}
}

func assertContains(t testing.TB, want, got string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Errorf("want substring %q in string %q, don't have it", want, got)
	}
}

func assertSameStruct(t testing.TB, want, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want event structure %+v, got %+v", want, got)
	}
}

func assertSameLength(t testing.TB, want, got int) {
	t.Helper()

	if want != got {
		t.Fatalf("want events length %d, got %d", want, got)
	}
}
