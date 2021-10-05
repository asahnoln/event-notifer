package pkg

import (
	"fmt"
	"strings"
	"time"
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

func Send(es []Event, sdr Sender, when EventType) error {
	return sdr.Send(generateMessage(es, when))
}

func generateMessage(es []Event, when EventType) string {
	message := &strings.Builder{}

	general := `
❗️ %s %s!

%s

`[1:]

	tmpl := `
%s
📍 Место: %s
👥 Требуются: %s
▶️ Начало: %s
⏹ Окончание: %s

`[1:]

	timeWord := time.Now().Format("02.01.2006")
	whenWord := "Сегодня"
	if when == Tomorrow {
		whenWord = "Завтра"
		timeWord = time.Now().AddDate(0, 0, 1).Format("02.01.2006")
	}
	what := "репетиция"
	if len(es) > 1 {
		what = "репетиции"
	}
	fmt.Fprintf(message, general, whenWord, what, timeWord)

	for _, e := range es {
		fmt.Fprintf(message, tmpl, e.What, e.Where, strings.Join(e.Who, ", "), e.Start, e.End)
	}

	return message.String()
}
