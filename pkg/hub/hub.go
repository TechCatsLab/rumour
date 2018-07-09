/*
 * Revision History:
 *     Initial: 2018/05/27        Tong Yuehong
 */

package hub

import (
	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/pkg/dispatcher"
	"github.com/TechCatsLab/rumour/pkg/manager"
	"github.com/TechCatsLab/rumour/pkg/queue"
	"github.com/TechCatsLab/rumour/pkg/log"
	"github.com/TechCatsLab/rumour/constants"
)

type hub struct {
	connectionManager rumour.ConnectionManager
	dispatcher        rumour.Dispatcher
	rumour.Queue
	shutdown chan struct{}
}

// NewHub - create a new Hub.
func NewHub() rumour.Hub {
	hub := &hub {
		shutdown: make(chan struct{}),
		Queue:    queue.NewChannelQueue(constants.HubQueueSize),
	}

	hub.connectionManager = manager.NewManager(hub)
	hub.dispatcher = dispatcher.NewDispatcher(hub)

	go hub.start()
	return hub
}

func (h *hub) start() {
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
func (h *hub) handleMessage() error {
	msg, err := h.Queue.Get()
	if err != nil {
		return err
	}

	if err = h.HubDispatcher().Dispatch(msg); err != nil {
		return err
	}

	return nil
}

// Dispatcher return the dispatcher.
func (h *hub) HubDispatcher() rumour.Dispatcher {
	return h.dispatcher
}

// ConnManager return connectionManager.
func (h *hub) ConnManager() rumour.ConnectionManager {
	return h.connectionManager
}

// Dispatch a message to the queue.
func (h *hub) Dispatch(message rumour.Message) error {
	h.Queue.Put(message)

	return nil
}
