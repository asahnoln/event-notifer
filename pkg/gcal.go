package pkg

import (
	"context"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GCalStore struct {
	calendarId string
	opts       []option.ClientOption
}

func NewGCalStore(calendarId string, opts ...option.ClientOption) *GCalStore {
	return &GCalStore{calendarId, opts}
}

func (s *GCalStore) Events(when EventType) ([]Event, error) {
	ctx := context.Background()

	srv, err := calendar.NewService(ctx, s.opts...)

	if err != nil {
		return nil, err
	}

	// TODO: Test time!
	es, err := srv.Events.List(s.calendarId).
		TimeMin(time.Now().AddDate(0, 0, 1).Format(time.RFC3339)).
		TimeMax(time.Now().AddDate(0, 0, 2).Format(time.RFC3339)).
		Do()

	if err != nil {
		return nil, err
	}

	var result []Event
	for _, e := range es.Items {
		var names []string
		for _, a := range e.Attendees {
			names = append(names, a.DisplayName)
		}

		result = append(result, Event{
			e.Summary,
			e.Location,
			names,
		})
	}

	return result, err
}
