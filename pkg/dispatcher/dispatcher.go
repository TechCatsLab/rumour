/*
 * Revision History:
 *     Initial: 2018/05/30        Tong Yuehong
 */

package dispatcher

import (
	"errors"
	"context"

	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/constants"
	"github.com/TechCatsLab/rumour/pkg/log"
	"github.com/TechCatsLab/scheduler"
)

const (
	ctxKeyMessage = "message"
)

var (
	ErrDispatch = errors.New("Can't dispatch the message.")
)

type dispatcher struct {
	pool *scheduler.Pool
	hub  rumour.Hub
}

// NewDispatcher new a dispatcher.
func NewDispatcher(hub rumour.Hub) rumour.Dispatcher {
	return &dispatcher{
		pool: scheduler.New(constants.PoolQueueSize, constants.PoolWorkerSize),
		hub:  hub,
	}
}

// Hub return hub.
func (d *dispatcher) Hub() rumour.Hub {
	return d.hub
}

func (d *dispatcher) Dispatch(message rumour.Message) error {
	ctx := context.WithValue(context.Background(), ctxKeyMessage, message)

	return d.pool.Schedule(scheduler.TaskFunc(d.dispatch), ctx)
}

func (d *dispatcher) dispatch(c context.Context) error {
	message := c.Value(ctxKeyMessage).(rumour.Message)
	userID := message.Target()

	conns, err := d.Hub().ConnManager().Query(userID)
	if err != nil {
		log.Error("[Dispatcher Dispatch] Query Connection err", log.Err(err))
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
