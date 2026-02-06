package push

import (
	"errors"

	"my-IMSystem/ws-gateway/internal/session"
)

type Dispatcher struct {
	manager *session.Manager
}

func NewDispatcher(manager *session.Manager) *Dispatcher {
	return &Dispatcher{manager: manager}
}

func (d *Dispatcher) DispatchToUser(userID int64, data []byte) error {
	if d.manager == nil {
		return errors.New("session manager is nil")
	}
	return d.manager.SendTo(userID, data)
}

func (d *Dispatcher) DispatchToDevice(userID int64, _ string, data []byte) error {
	return d.DispatchToUser(userID, data)
}

func (d *Dispatcher) Broadcast(data []byte) {
	if d.manager == nil {
		return
	}
	d.manager.Broadcast(data)
}
