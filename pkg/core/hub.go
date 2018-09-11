/*
 * Revision History:
 *     Initial: 2018/05/27        Tong Yuehong
 */

package core

import (
	"context"
	"errors"

	"github.com/TechCatsLab/scheduler"

	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/message"
	"github.com/TechCatsLab/rumour/pkg/generator"
	"github.com/TechCatsLab/rumour/pkg/queue"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
)

const (
	ctxKeyMessage = "single_message"
)

var (
	ErrDispatch = errors.New("Can't dispatch the message.")
)

// Hub represents a core.
type Hub struct {
	ConnectionManager *ConnectionManager
	ChannelManager    *Channels

	store *mysql.StoreServiceProvider
	pool  *scheduler.Pool

	rumour.Queue
	generator *generator.Generator
	shutdown  chan struct{}
}

// NewHub - create a new Hub.
func newHub(c *Config) *Hub {
	hub := &Hub{
		ChannelManager:    NewChannelManager(),
		ConnectionManager: NewConnectionManager(),
		pool:              scheduler.New(c.DispatcherQueueSize, c.DispatcherWorkers),
		shutdown:          make(chan struct{}),
		Queue:             queue.NewChannelQueue(c.IncomingMessageQueueSize),
	}

	// TODO: Load from configuration or database.
	hub.generator = generator.New(1000, hub.shutdown)
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
			log.Error(err)
		}
	}
}

func (h *Hub) handleMessage() error {
	msg, err := h.Queue.Get()
	if err != nil {
		return err
	}

	switch msg.Type {
	case message.MessageTypeText:
		h.Dispatch(msg)
	case message.MessageTypeChanText:
		h.ChannelMessage(msg)
	}

	return nil
}

//Generator return generator.
func (h *Hub) Generator() *generator.Generator {
	return h.generator
}

// Dispatch a message to the queue.
func (h *Hub) Put(message *rumour.Message) error {
	return h.Queue.Put(message)
}

// JoinChannel add connections to the channel.
func (hub *Hub) JoinChannel(userID string, chanID uint32) error {
	conns, err := hub.ConnectionManager.Query(userID)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		err := hub.ChannelManager.Add(chanID, conn)
		if err != nil {
			return err
		}
	}

	return nil
}

// ChannelMessage dispatch channel message.
func (hub *Hub) ChannelMessage(message *rumour.Message) error {
	return hub.ChannelManager.Dispatch(message)
}

// Dispatch dispatch the single message.
func (hub *Hub) Dispatch(message *rumour.Message) error {
	mysql.StoreService.Store().SingleMessage().Insert(
		message.Content["id"].(uint64),
		message.From,
		message.To,
		uint32(message.Type),
		message.Content["message"].(string))
	ctx := context.WithValue(context.Background(), ctxKeyMessage, message)

	return hub.pool.Schedule(ctx, scheduler.TaskFunc(hub.dispatch))
}

func (hub *Hub) dispatch(ctx context.Context) error {
	message := ctx.Value(ctxKeyMessage).(*rumour.Message)
	userID := message.To

	conns, err := hub.ConnectionManager.Query(userID)
	if err != nil {
		log.Error(err)
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
