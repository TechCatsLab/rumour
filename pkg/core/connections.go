/*
 * Revision History:
 *     Initial: 2018/05/31        Tong Yuehong
 */

package core

import (
	"errors"
	"sync"

	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour"
)

var ErrConnNotExist = errors.New("connection not exists")

type ConnectionManager struct {
	conns map[string][]rumour.Connection
	mux   sync.RWMutex
}

// NewManager create a new connection manager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		conns: make(map[string][]rumour.Connection),
	}
}

// Add a connection.
func (m *ConnectionManager) Add(connection rumour.Connection) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	userID, err := connection.Identify()
	if err != nil {
		log.Warn(err)
		connection.Stop()
		return nil
	}

	if _, exists := m.conns[userID]; !exists {
		m.conns[userID] = []rumour.Connection{}
	}
	m.conns[userID] = append(m.conns[userID], connection)

	return nil
}

// Remove a connection.
func (m *ConnectionManager) Remove(connection rumour.Connection) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	userID, err := connection.Identify()
	if err != nil {
		connection.Stop()
		return nil
	}

	conns, exists := m.conns[userID]
	if !exists {
		return nil
	}

	for i, v := range conns {
		if v == connection {
			conns = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	m.conns[userID] = conns
	return nil
}

// Query someone's connection.
func (m *ConnectionManager) Query(id string) ([]rumour.Connection, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	conns, exists := m.conns[id]
	if !exists {
		log.Error("[Manager Query] Connection not exists")
		return nil, ErrConnNotExist
	}

	return conns, nil
}
