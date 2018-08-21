/*
 * Revision History:
 *     Initial: 2018/07/29        Tong Yuehong
 */

package core

import (
	"errors"
	"strconv"
	"sync"
	"context"

	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
	"github.com/TechCatsLab/scheduler"
)

var (
	ErrChannelNotExist = errors.New("channel manager: channel not exists")
)

type Channels struct {
	muCh     sync.RWMutex
	channels map[uint32]*Channel
	pool     *scheduler.Pool
}

// NewManager creates a channel manager.
func NewChannelManager() *Channels {
	return &Channels{
		channels: make(map[uint32]*Channel),
		pool:     scheduler.New(16, 10),
	}
}

// Load all channels from database.
func (chans *Channels) Load() error {
	channels, err := mysql.StoreService.Store().Channel().QueryExist()
	if err != nil {
		return err
	}

	chans.muCh.Lock()
	for _, channel := range channels {
		c := NewChan(channel.Id)
		chans.channels[channel.Id] = c
	}
	chans.muCh.Unlock()

	return nil
}

// Add a channel.
func (chans *Channels) Add(chanID uint32, conn rumour.Connection) error {
	chans.muCh.Lock()
	defer chans.muCh.Unlock()

	if _, exists := chans.channels[chanID]; !exists {
		chans.channels[chanID] = NewChan(chanID)
	}

	err := chans.channels[chanID].Add(conn)
	if err != nil {
		return err
	}

	return nil
}

// Remove a channel.
func (chans *Channels) Remove(chanID uint32) error {
	chans.muCh.Lock()
	if _, exists := chans.channels[chanID]; !exists {
		chans.muCh.Unlock()
		return nil
	}

	delete(chans.channels, chanID)
	chans.muCh.Unlock()

	// TODO: Tx
	err := mysql.StoreService.Store().Channel().Disable(chanID)
	if err != nil {
		return err
	}

	err = mysql.StoreService.Store().ChannelUser().Remove(chanID)
	if err != nil {
		return err
	}

	return nil
}

// Query
func (chans *Channels) Query(id string) (*Channel, error) {
	chanID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	chans.muCh.RLock()
	defer chans.muCh.RUnlock()

	channel, exists := chans.channels[uint32(chanID)]
	if !exists {
		return nil, ErrChannelNotExist
	}

	return channel, nil
}

func (chans *Channels) Dispatch(message *rumour.Message) error {
	_, err := mysql.StoreService.Store().ChannelMessage().Insert(
		message.Content["id"].(uint64),
		message.From,
		message.To,
		uint8(message.Content["kind"].(float64)),
		message.Content["message"].(string))

	if err != nil {
		return err
	}


	ctx := context.WithValue(context.Background(), ctxKeyMessage, message)

	return chans.pool.Schedule(ctx, scheduler.TaskFunc(chans.dispatch))
}

// Channel message dispatch task.
func (chans Channels) dispatch(ctx context.Context) error {
	message := ctx.Value(ctxKeyMessage).(*rumour.Message)
	to := message.To

	channel, err := chans.Query(to)
	if err != nil {
		return err
	}

	return channel.Send(message)
}
