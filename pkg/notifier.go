package pkg

import (
	"fmt"
	"strings"
)

type EventType int

const (
	Today EventType = iota
	Tomorrow
)

type Event struct {
	What, Where, Start, End string
	Who                     []string
}

type Store interface {
	Events(when EventType) ([]Event, error)
}

type Sender interface {
	Send(message string) error
}

func TomorrowEvents(store Store) ([]Event, error) {
	return store.Events(Tomorrow)
}

func TodayEvents(store Store) ([]Event, error) {
	return store.Events(Today)
}

func Send(es []Event, sdr Sender) error {
	return sdr.Send(generateMessage(es))
}

func generateMessage(es []Event) string {
	message := &strings.Builder{}

	general := `
❗️ Завтра %s!

08.11.2019

`[1:]

	tmpl := `
%s
📍 Место: %s
👥 Требуются: %s
▶️ Начало: %s
⏹ Окончание: %s

`[1:]

	what := "репетиция"
	if len(es) > 1 {
		what = "репетиции"
	}
	fmt.Fprintf(message, general, what)
	for _, e := range es {
		fmt.Fprintf(message, tmpl, e.What, e.Where, strings.Join(e.Who, ", "), e.Start, e.End)
	}

	return message.String()
}
