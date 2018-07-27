/*
 * Revision History:
 *     Initial: 2018/07/26        Tong Yuehong
 */

package store

import (
	"time"
)

type (
	singleMessageServiceProvider struct {
		store *Store
	}

	singleMessage struct {
		Id        uint64
		SourceID  uint64
		TargetID  uint64
		Kind      uint
		Content   string
		CreatedAt time.Time
		Arrived   bool
	}
)

// Create single_message table.
func (sm *singleMessageServiceProvider) Create() error {
	_, err := sm.store.DB.Exec(sqlStmts[SQLSingleMessageCreateTable])
	return err
}

// Insert single message.
func (sm *singleMessageServiceProvider) Insert(id, from, to uint64, kind uint, content *string) error {
	result, err := sm.store.DB.Exec(sqlStmts[SQLSingleMessageInsert], id, from, to, kind, content)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidInsert
	}

	return nil
}

// QueryRecord query chat record by TargetID and sourceID.
func (sm *singleMessageServiceProvider) QueryRecord(from, to uint64) ([]*singleMessage, error) {
	var (
		id        uint64
		sourceID  uint64
		targetID  uint64
		kind      uint
		content   string
		createdAt time.Time
		record    []*singleMessage
	)

	rows, err := sm.store.DB.Query(sqlStmts[SQLSingleMessageQueryRecord], from, to, to, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &sourceID, &targetID, &kind, &content, &createdAt); err != nil {
			return nil, err
		}

		single := &singleMessage{
			Id:        id,
			SourceID:  sourceID,
			TargetID:  targetID,
			Kind:      kind,
			Content:   content,
			CreatedAt: createdAt,
		}
		record = append(record, single)
	}

	return record, nil
}

// QueryUnread query unread message by TargetID.
func (sm *singleMessageServiceProvider) QueryUnreadByUserID(from, to uint64) ([]*singleMessage, error) {
	var (
		id        uint64
		sourceID  uint64
		targetID  uint64
		kind      uint
		content   string
		createdAt time.Time
		unread    []*singleMessage
	)

	rows, err := sm.store.DB.Query(sqlStmts[SQLSingleMessageUnreadByUserID], from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &sourceID, &targetID, &kind, &content, &createdAt); err != nil {
			return nil, err
		}

		single := &singleMessage{
			Id:        id,
			SourceID:  sourceID,
			TargetID:  targetID,
			Kind:      kind,
			Content:   content,
			CreatedAt: createdAt,
		}

		unread = append(unread, single)
	}

	return unread, nil
}

// QueryUnread query unread message by TargetID.
func (sm *singleMessageServiceProvider) QueryUnreadByMsgID(from, to, msgID uint64) ([]*singleMessage, error) {
	var (
		id        uint64
		sourceID  uint64
		targetID  uint64
		kind      uint
		content   string
		createdAt time.Time
		unread    []*singleMessage
	)

	rows, err := sm.store.DB.Query(sqlStmts[SQLSingleMessageUnreadByMsgID], to, msgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &sourceID, &targetID, &kind, &content, &createdAt); err != nil {
			return nil, err
		}

		single := &singleMessage{
			Id:        id,
			SourceID:  sourceID,
			TargetID:  targetID,
			Kind:      kind,
			Content:   content,
			CreatedAt: createdAt,
		}

		unread = append(unread, single)
	}

	return unread, nil
}
