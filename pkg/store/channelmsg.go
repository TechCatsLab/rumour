/*
 * Revision History:
 *     Initial: 2018/07/26        Tong Yuehong
 */

package store

import "time"

type (
	channelMessageServiceProvider struct {
		store *Store
	}

	channelMessage struct {
		Id        uint64
		SourceID  uint64
		ChannelID uint
		Kind      uint8
		Content   string
		CreatedAt time.Time
	}
)

// Create table channel_message.
func (cm *channelMessageServiceProvider) Create() error {
	_, err := cm.store.DB.Exec(sqlStmts[SQLChannelMessageCreateTable])

	return err
}

// Insert a raw.
func (cm *channelMessageServiceProvider) Insert(id, from uint64, to uint, kind uint, content *string) error {
	result, err := cm.store.DB.Exec(sqlStmts[SQLChannelMessageInsert], id, from, to, kind, content)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidInsert
	}

	return nil
}

// Unread return the message by messageID.
func (cm *channelUserServiceProvider) Unread(chanID uint, msgID uint64) ([]*channelMessage, error) {
	var (
		id        uint64
		sourceID  uint64
		kind      uint8
		content   string
		createdAt time.Time
		unread    []*channelMessage
	)

	rows, err := cm.store.DB.Query(sqlStmts[SQLChannelMessageUnread], chanID, msgID)
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
