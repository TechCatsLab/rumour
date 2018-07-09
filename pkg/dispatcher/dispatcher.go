/*
 * Revision History:
 *     Initial: 2018/05/30        Tong Yuehong
 */

package dispatcher

import (
	"errors"

	"github.com/TechCatsLab/rumour"
)

var ErrDispatch = errors.New("Can't dispatch the message")

type dispatcher struct {
	hub rumour.Hub
}

// NewQueue - new a messageQueue.
func NewDispatcher(hub rumour.Hub) rumour.Dispatcher {
	return &dispatcher{hub}
}

// Hub -
func (d *dispatcher) Hub() rumour.Hub {
	return d.hub
}

// Dispatch -
func (d *dispatcher) Dispatch(message rumour.Message) error {
	userID := message.Target()

	conns, err := d.Hub().ConnManager().Query(userID)
	if err != nil {
		return err
	}

	succeed := false
	for _, conn := range conns {
		err = conn.Send(message)
		if err != nil {
			continue
		}

		succeed = true
	}

	if succeed {
		return nil
	}

	for _, conn := range conns {
		conn.Stop()
	}

	return ErrDispatch
}
