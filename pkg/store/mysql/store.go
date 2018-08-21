/*
 * Revision History:
 *     Initial: 2018/07/25        Tong Yuehong
 */

package mysql

import (
	"database/sql"
)

type StoreServiceProvider struct {
	s *store
}

var (
	StoreService *StoreServiceProvider
)

func NewStore(db *sql.DB) error {
	s, err := newStore(db)
	if err != nil {
		return err
	}

	StoreService = &StoreServiceProvider{
		s:s,
	}

	return nil
}

func (sp *StoreServiceProvider) Store() *store {
	return sp.s
}

type store struct {
	db *sql.DB

	channel        *channelServiceProvider
	singleMessage  *singleMessageServiceProvider
	channelMessage *channelMessageServiceProvider
	channelUser    *channelUserServiceProvider
}

// newStore new a store.
func newStore(db *sql.DB) (*store, error) {
	s := &store{
		db: db,
	}

	s.channel = &channelServiceProvider{
		store: s,
	}

	s.singleMessage = &singleMessageServiceProvider{
		store: s,
	}

	s.channelMessage = &channelMessageServiceProvider{
		store: s,
	}

	s.channelUser = &channelUserServiceProvider{
		store: s,
	}

	return s, nil
}

// Channel returns channelServiceProvider.
func (s *store) Channel() *channelServiceProvider {
	return s.channel
}

// SingleMessage returns singleMessageProvider.
func (s *store) SingleMessage() *singleMessageServiceProvider {
	return s.singleMessage
}

// ChannelMessage returns channelMessageProvider.
func (s *store) ChannelMessage() *channelMessageServiceProvider {
	return s.channelMessage
}

// ChannelUser returns channelUserProvider.
func (s *store) ChannelUser() *channelUserServiceProvider {
	return s.channelUser
}
