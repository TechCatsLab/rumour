/*
 * Revision History:
 *     Initial: 2018/07/25        Tong Yuehong
 */

package store

import (
	"database/sql"
)

type Store struct {
	DB             *sql.DB

	channel        *channelServiceProvider
	singleMessage  *singleMessageServiceProvider
	channelMessage *channelMessageServiceProvider
	channelUser    *channelUserServiceProvider
}

// NewStore new a store.
func NewStore(dataSource string) (*Store, error) {
	chatDB, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	return &Store{
		DB: chatDB,
	}, nil
}

// Channel returns channelServiceProvider.
func (s *Store) Channel() *channelServiceProvider {
	if s.channel == nil {
		s.channel = &channelServiceProvider{
			store: s,
		}
	}

	return s.channel
}

// SingleMessage returns singleMessageProvider.
func (s *Store) SingleMessage() *singleMessageServiceProvider {
	if s.singleMessage == nil {
		s.singleMessage = &singleMessageServiceProvider{
			store: s,
		}
	}

	return s.singleMessage
}

// ChannelMessage returns channelMessageProvider.
func (s *Store) ChannelMessage() *channelMessageServiceProvider {
	if s.channelMessage == nil {
		s.channelMessage = &channelMessageServiceProvider{
			store: s,
		}
	}

	return s.channelMessage
}

// ChannelUser returns channelUserProvider.
func (s *Store) ChannelUser() *channelUserServiceProvider {
	if s.channelUser == nil {
		s.channelUser = &channelUserServiceProvider{
			store: s,
		}
	}

	return s.channelUser
}
