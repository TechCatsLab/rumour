/*
 * Revision History:
 *     Initial: 2018/07/25        Tong Yuehong
 */

package core

import (
	"sync"

	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
)

// Channel represents a channel.
type Channel struct {
	id          uint32
	muConn      sync.RWMutex
	connections []rumour.Connection
}

// NewChan create a new channel.
func NewChan(id uint32) *Channel {
	channel :=  &Channel{
		id:          id,
		connections: make([]rumour.Connection, 0),
	}

	return channel
}

// Add a connection to a channel.
func (ch *Channel) Add(conn rumour.Connection) error {
	ch.muConn.Lock()
	defer ch.muConn.Unlock()

	ch.connections = append(ch.connections, conn)
	return nil
}

// Remove a connection from a channel.
func (ch *Channel) Remove(conn rumour.Connection) error {
	ch.muConn.Lock()
	defer ch.muConn.Unlock()

	for i, connection := range ch.connections {
		if connection == conn {
			ch.connections = append(ch.connections[:i], ch.connections[i+1:]...)
			break
		}
	}

	return nil
}

// Send a message to all connections in channel.
func (ch *Channel) Send(message *rumour.Message) error {
	ch.muConn.RLock()
	defer ch.muConn.RUnlock()

	for _, conn := range ch.connections {
		err := conn.Send(message)
		if err != nil {
			continue
		}

		userID, err := conn.Identify()
		if err != nil {
			continue
		}

		err = mysql.StoreService.Store().ChannelUser().UpdateMsgID(ch.id, message.Content["id"].(uint64), userID)
		if err != nil {
			continue
		}
	}

	return nil
}
