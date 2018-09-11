/*
 * Revision History:
 *     Initial: 2018/07/26        Tong Yuehong
 */

package mysql

import (
	"strconv"
	"time"
	"errors"
)

type (
	singleMessageServiceProvider struct {
		store *store
	}

	singleMessage struct {
		Id        uint64
		SourceID  uint64
		TargetID  uint64
		Kind      uint32
		Content   string
		CreatedAt time.Time
		Arrived   bool
	}
)

const (
	sqlSingleMessageCreateTable                    = iota
	sqlSingleMessageInsert
	sqlSingleMessageQueryRecord
	sqlSingleMessageUnreadByUserID
	sqlSingleMessageUnreadByMsgID
)
var (
	sqlSingle_Message = []string{
		`CREATE TABLE IF NOT EXISTS chat.single_message (
			id 			BIGINT UNSIGNED NOT NULL,
			source_id 	BIGINT UNSIGNED NOT NULL DEFAULT 0,
			target_id 	BIGINT UNSIGNED NOT NULL DEFAULT 0,
			kind 		INTEGER UNSIGNED NOT NULL DEFAULT 0,
			content     TEXT NOT NULL,
			created_at 	DATETIME NOT NULL DEFAULT current_timestamp,
			arrived 	BOOL NOT NULL DEFAULT false,
			INDEX (id),
			INDEX (source_id, target_id),
			PRIMARY KEY (Id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO chat.single_message(id, source_id,target_id,kind,content) VALUES (?,?,?,?,?)`,
		`SELECT id,source_id,target_id,kind,content,created_at FROM chat.single_message WHERE (source_id = ? AND target_id = ?) OR (source_id = ? AND target_id = ?) LOCK IN SHARE MODE`,
		`SELECT id,source_id,target_id,kind,content,created_at FROM chat.single_message WHERE source_id = ? AND target_id = ? AND arrived = false LOCK IN SHARE MODE`,
		`SELECT id,source_id,target_id,kind,content,created_at FROM chat.single_message WHERE id > ? AND source_id = ? AND target_id = ? AND arrived = false LOCK IN SHARE MODE`,
	}
)

var (
	errInvalidSMInsert = errors.New("insert single_message: insert single message affected 0 row")
)

// Create single_message table.
func (sm *singleMessageServiceProvider) Create() error {
	_, err := sm.store.db.Exec(sqlSingle_Message[sqlSingleMessageCreateTable])
	return err
}

// Insert single message.
func (sm *singleMessageServiceProvider) Insert(id uint64, from, to string, kind uint32, content string) (uint64, error) {
	source, err := strconv.ParseUint(from, 10, 64)
	if err != nil {
		return 0, err
	}

	target, err := strconv.ParseUint(to, 10, 64)
	if err != nil {
		return 0, err
	}

	result, err := sm.store.db.Exec(sqlSingle_Message[sqlSingleMessageInsert], id, source, target, kind, content)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvalidSMInsert
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ID), nil
}

// QueryRecord query chat record by TargetID and sourceID.
func (sm *singleMessageServiceProvider) QueryRecord(from, to string) ([]*singleMessage, error) {
	var (
		id        uint64
		sourceID  uint64
		targetID  uint64
		kind      uint32
		content   string
		createdAt time.Time
		record    []*singleMessage
	)

	source, err := strconv.ParseUint(from, 10, 64)
	if err != nil {
		return nil, err
	}

	target, err := strconv.ParseUint(to, 10, 64)
	if err != nil {
		return nil, err
	}

	rows, err := sm.store.db.Query(sqlSingle_Message[sqlSingleMessageQueryRecord], source, target, to, from)
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
func (sm *singleMessageServiceProvider) QueryUnreadByUserID(from, to string) ([]*singleMessage, error) {
	var (
		id        uint64
		sourceID  uint64
		targetID  uint64
		kind      uint32
		content   string
		createdAt time.Time
		unread    []*singleMessage
	)

	source, err := strconv.ParseUint(from, 10, 64)
	if err != nil {
		return nil, err
	}

	target, err := strconv.ParseUint(to, 10, 64)
	if err != nil {
		return nil, err
	}

	rows, err := sm.store.db.Query(sqlSingle_Message[sqlSingleMessageUnreadByUserID], source, target)
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
func (sm *singleMessageServiceProvider) QueryUnreadByMsgID(from, to string, msgID uint64) ([]*singleMessage, error) {
	var (
		id        uint64
		sourceID  uint64
		targetID  uint64
		kind      uint32
		content   string
		createdAt time.Time
		unread    []*singleMessage
	)

	source, err := strconv.ParseUint(from, 10, 64)
	if err != nil {
		return nil, err
	}

	target, err := strconv.ParseUint(to, 10, 64)
	if err != nil {
		return nil, err
	}

	rows, err := sm.store.db.Query(sqlSingle_Message[sqlSingleMessageUnreadByMsgID], msgID, source, target)
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
