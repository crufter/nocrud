package events

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	gt "github.com/opesun/gotrigga"
	"fmt"
)

var conn_never = fmt.Errorf("Connection to server was never established.")

type Events struct {
	conn	*gt.Connection
}

func New (conn *gt.Connection) *Events {
	return &Events{
		conn,
	}
}

type Event struct {
	room	*gt.Room
	name	string
}

func (e *Events) Select(name string) iface.Event {
	if e.conn == nil {
		return &Event{
			nil,
			name,
		}
	}
	room := e.conn.Room(name)
	return &Event{
		room,
		name,
	}
}

func (e *Event) Publish(msg []byte) error {
	if e.room == nil {
		return conn_never
	}
	return e.room.Publish(string(msg))
}

func (e *Event) Subscribe() error {
	if e.room == nil {
		return conn_never
	}
	return e.room.Subscribe()
}

func (e *Event) Unsubscribe() error {
	if e.room == nil {
		return conn_never
	}
	return e.room.Unsubscribe()
}

func (e *Event) Read() ([]byte, error) {
	if e.room == nil {
		return nil, conn_never
	}
	return e.room.Read()
}
