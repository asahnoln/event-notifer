package pkg_test

import (
	"testing"
	"time"

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
	err := pkg.Send(es, sdr, pkg.Tomorrow)
	assertNoError(t, err, "unexpected error while sending message: %v")

	want := es[0]
	assertContains(t, want.What, sdr.result)
	assertContains(t, want.Where, sdr.result)
	for _, p := range want.Who {
		assertContains(t, p, sdr.result)
	}
}

func TestDayWordInMessage(t *testing.T) {
	tests := []struct {
		eventType pkg.EventType
		want      string
	}{
		{pkg.Tomorrow, "Завтра"},
		{pkg.Today, "Сегодня"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			es, _ := pkg.TomorrowEvents(db)
			sdr := &stubSender{}
			err := pkg.Send(es, sdr, tt.eventType)
			assertNoError(t, err, "unexpected error while sending message: %v")
			assertContains(t, tt.want, sdr.result)
		})
	}
}

func TestProperDateInMessageToday(t *testing.T) {
	tests := []struct {
		eventsFunc func(pkg.Store) ([]pkg.Event, error)
		eventType  pkg.EventType
		date       string
	}{
		{pkg.TodayEvents, pkg.Today, time.Now().Format("02.01.2006")},
		{pkg.TomorrowEvents, pkg.Tomorrow, time.Now().AddDate(0, 0, 1).Format("02.01.2006")},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			es, _ := tt.eventsFunc(db)
			sdr := &stubSender{}
			err := pkg.Send(es, sdr, tt.eventType)
			assertNoError(t, err, "unexpected error while sending message: %v")

			assertContains(t, tt.date, sdr.result)
		})
	}
}

func TestSendErrorWhenEventsEmpty(t *testing.T) {
	sdr := &stubSender{}
	err := pkg.Send([]pkg.Event{}, sdr, pkg.Tomorrow)

	assertError(t, err, "want error because there were 0 events sent, but got nil")
}
