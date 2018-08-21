/*
 * Revision History:
 *     Initial: 2018/07/25        Tong Yuehong
 */

package core

import (
	"sync"

	"github.com/TechCatsLab/rumour"
)

type Channel struct {
	id          uint32
	muConn      sync.RWMutex
	connections []rumour.Connection
}

func NewChan(id uint32) *Channel {
	channel :=  &Channel{
		id:          id,
		connections: make([]rumour.Connection, 0),
	}

	return channel
}

func (ch *Channel) Add(conn rumour.Connection) error {
	ch.muConn.Lock()
	defer ch.muConn.Unlock()

	ch.connections = append(ch.connections, conn)
	return nil
}

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

func (ch *Channel) Send(message *rumour.Message) error {
	ch.muConn.RLock()
	defer ch.muConn.RUnlock()

	for _, conn := range ch.connections {
		err := conn.Send(message)
		if err != nil {
			continue
		}
	}

	return nil
}
