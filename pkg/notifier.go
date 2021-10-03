package pkg

type EventType int

const (
	Today EventType = iota
	Tomorrow
)

type Event struct {
	What, Where string
	Who         []string
}

type Store interface {
	Events(when EventType) []Event
}

func TomorrowEvents(store Store) []Event {
	return store.Events(Tomorrow)
}

func TodayEvents(store Store) []Event {
	return store.Events(Today)
}
