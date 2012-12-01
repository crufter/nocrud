package hooks

import (
	iface "github.com/opesun/nocrud/frame/interfaces"
)

// Used to call subscribed hooks.
type Hooks struct {
}

type Hook struct {
}

func New() *Hooks {
	return &Hooks{
	}
}

func (e *Hooks) Initer(initer func(iface.Instance) error) {
}

func (e *Hooks) Select(Hookname string) iface.Hook {
	return &Hook{
	}
}

type subscriber struct {
}

func (s *subscriber) Name() string {
	return ""
}

func (s *subscriber) Method() string {
	return ""
}

func (e *Hook) Subscribers() []iface.Subscriber {
	return []iface.Subscriber{}
}

type InstanceCacher struct {
	iface.Module
}

func (m InstanceCacher) Instance() iface.Instance {
	return nil
}

func (e *Hooks) Module(modname string) iface.Module {
	return &InstanceCacher{
	}
}

func (e *Hook) HasSubscribers() bool {
	return false
}

func (e *Hook) SubscriberCount() int {
	return 0
}

func (e *Hook) Fire(params ...interface{}) {
}

func (e *Hook) Iterate(stopfunc interface{}, params ...interface{}) {
}
