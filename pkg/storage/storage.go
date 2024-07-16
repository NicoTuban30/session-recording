package storage

type Event any

type Storage interface {
	Events(session string) ([]Event, error)

	AppendEvents(session string, events ...Event) error
	DeleteSession(sessionID string) error
}
