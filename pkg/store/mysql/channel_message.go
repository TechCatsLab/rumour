/*
 * Revision History:
 *     Initial: 2018/07/26        Tong Yuehong
 */

package mysql

import (
	"strconv"
	"time"
)

type (
	channelMessageServiceProvider struct {
		store *store
	}

	channelMessage struct {
		Id        uint64
		SourceID  uint64
		ChannelID uint32
		Kind      uint8
		Content   string
		CreatedAt time.Time
	}
)

const (
	sqlChannelMessageCreateTable   = iota
	sqlChannelMessageInsert
	sqlChannelMessageUnread
)

var (
	sqlChannel_Message = []string{
		`CREATE TABLE IF NOT EXISTS chat.channel_message (
			id 			BIGINT UNSIGNED NOT NULL,
			channel_id	INTEGER UNSIGNED NOT NULL DEFAULT 0,
			source_id 	BIGINT UNSIGNED NOT NULL DEFAULT 0,
			kind 		INTEGER UNSIGNED NOT NULL DEFAULT 0,
			content		TEXT NOT NULL,
			created_at 	DATETIME NOT NULL DEFAULT current_timestamp,
			INDEX (channel_id),
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO chat.channel_message(Id,source_id,channel_id,kind,content) VALUES (?,?,?,?,?)`,
		`SELECT id,source_id,kind,content,created_at FROM chat.channel_message WHERE channel_id = ? AND id > ? LOCK IN SHARE MODE`,
	}
)

// Create table channel_message.
func (cm *channelMessageServiceProvider) Create() error {
	_, err := cm.store.db.Exec(sqlChannel_Message[sqlChannelMessageCreateTable])

	return err
}

// Insert a raw.
func (cm *channelMessageServiceProvider) Insert(id uint64, from, to string, kind uint8, content string) (uint64, error) {
	sourceID, err := strconv.ParseUint(from, 10, 64)
	if err != nil {
		return 0, err
	}

	targetID, err := strconv.ParseUint(to, 10, 64)
	if err != nil {
		return 0, err
	}

	result, err := cm.store.db.Exec(sqlChannel_Message[sqlChannelMessageInsert], id, sourceID, targetID, kind, content)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvalidInsert
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ID), nil
}

// Unread return the message by messageID.
func (cm *channelMessageServiceProvider) Unread(chanID uint32, msgID uint64) ([]*channelMessage, error) {
	var (
		id        uint64
		sourceID  uint64
		kind      uint8
		content   string
		createdAt time.Time
		unread    []*channelMessage
	)

	rows, err := cm.store.db.Query(sqlChannel_Message[sqlChannelMessageUnread], chanID, msgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &sourceID, &kind, &content, &createdAt); err != nil {
			return nil, err
		}

		message := &channelMessage{
			Id:        id,
			SourceID:  sourceID,
			Kind:      kind,
			Content:   content,
			CreatedAt: createdAt,
		}

		unread = append(unread, message)
	}

	return unread, nil
}
