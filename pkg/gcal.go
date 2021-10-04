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

	var (
		min, max time.Time
	)
	switch when {
	case Today:
		min = time.Now()
		max = min.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	default:
		min = time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		max = min.AddDate(0, 0, 1)
	}

	es, err := srv.Events.List(s.calendarId).
		TimeMin(min.Format(time.RFC3339)).
		TimeMax(max.Format(time.RFC3339)).
		Do()

	if err != nil {
		return nil, err
	}

	// TODO: Allocate the same length immediately as es.Items
	var result []Event
	for _, e := range es.Items {
		var names []string
		for _, a := range e.Attendees {
			names = append(names, a.DisplayName)
		}

		start, err := time.Parse(time.RFC3339, e.Start.DateTime)
		if err != nil {
			return result, err
		}
		end, err := time.Parse(time.RFC3339, e.End.DateTime)
		if err != nil {
			return result, err
		}

		result = append(result, Event{
			e.Summary,
			e.Location,
			start.Format("02.01.2006 15:04"),
			end.Format("02.01.2006 15:04"),
			names,
		})
	}

	return result, err
}
