/*
 * Revision History:
 *     Initial: 2018/05/27        Tong Yuehong
 */

package hub

import (
	"context"
	"errors"

	"github.com/TechCatsLab/scheduler"

	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/pkg/log"
	"github.com/TechCatsLab/rumour/pkg/queue"
)

const (
	ctxKeyMessage = "single_message"
)

var (
	ErrDispatch = errors.New("Can't dispatch the message.")
)

type Hub struct {
	ConnectionManager *ConnectionManager
	pool              *scheduler.Pool
	rumour.Queue
	shutdown chan struct{}
}

// NewHub - create a new Hub.
func newHub(c *Config) *Hub {
	hub := &Hub{
		ConnectionManager: NewConnectionManager(),
		pool:              scheduler.New(c.DispatcherQueueSize, c.DispatcherWorkers),
		shutdown:          make(chan struct{}),
		Queue:             queue.NewChannelQueue(c.IncomingMessageQueueSize),
	}

	hub.ConnectionManager = NewConnectionManager()

	go hub.start()
	return hub
}

func (h *Hub) start() {
	for {
		select {
		case <-h.shutdown:
			break
		default:
		}

		if err := h.handleMessage(); err != nil {
			log.Error("[Hub] handleMessage error", log.Err(err))
		}
	}
}

//
func (h *Hub) handleMessage() error {
	msg, err := h.Queue.Get()
	if err != nil {
		return err
	}

	h.Dispatch(msg)

	return nil
}

// Dispatch a message to the queue.
func (h *Hub) Put(message rumour.Message) error {
	h.Queue.Put(message)

	return nil
}


func (hub *Hub) Dispatch(message rumour.Message) error {
	ctx := context.WithValue(context.Background(), ctxKeyMessage, message)

	return hub.pool.Schedule(scheduler.TaskFunc(hub.dispatch), ctx)
}

func (hub *Hub) dispatch(c context.Context) error {
	message := c.Value(ctxKeyMessage).(rumour.Message)
	userID := message.Target()

	conns, err := hub.ConnectionManager.Query(userID)
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
