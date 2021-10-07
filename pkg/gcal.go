package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GCalStore struct {
	calendarId  string
	mailsReader io.Reader
	opts        []option.ClientOption
}

func NewGCalStore(calendarId string, mails io.Reader, opts ...option.ClientOption) *GCalStore {
	return &GCalStore{calendarId, mails, opts}
}

func (s *GCalStore) Events(when EventType) ([]Event, error) {
	ctx := context.Background()
	srv, err := calendar.NewService(ctx, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("gcal: service error: %w", err)
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
		ShowDeleted(false).
		SingleEvents(true).
		OrderBy("startTime").
		Do()

	if err != nil {
		return nil, fmt.Errorf("gcal: events error: %w", err)
	}

	mailsSource, err := io.ReadAll(s.mailsReader)
	if err != nil {
		return nil, fmt.Errorf("gcal: mails file reading error: %w", err)
	}

	var result []Event
	for _, e := range es.Items {
		mails := make([]string, len(e.Attendees))
		for i, a := range e.Attendees {
			mails[i] = a.Email
		}

		// TODO: Quick Dirty fix for rereading reader
		names, err := MailsToNames(mails, bytes.NewReader(mailsSource))
		if err != nil {
			return result, fmt.Errorf("gcal: mails converting error: %w", err)
		}

		start, err := time.Parse(time.RFC3339, e.Start.DateTime)
		if err != nil {
			return result, fmt.Errorf("gcal: start time parse error: %w", err)
		}
		end, err := time.Parse(time.RFC3339, e.End.DateTime)
		if err != nil {
			return result, fmt.Errorf("gcal: end time parse error: %w", err)
		}

		result = append(result, Event{
			e.Summary,
			e.Location,
			start.Format("02.01.2006 15:04"),
			end.Format("02.01.2006 15:04"),
			names,
		})
	}

	return result, nil
}
