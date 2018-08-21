/*
 * Revision History:
 *     Initial: 2018/07/26        Tong Yuehong
 */

package mysql

import (
	"strconv"
)

type (
	channelUserServiceProvider struct {
		store *store
	}

	channelUser struct {
		Id            uint64
		ChannelID     uint32
		UserID        uint64
		Role          uint8
		LastMessageID uint64
	}
)

const (
	sqlChannelUserCreateTable      = iota
	sqlChannelUserInsert
	sqlChannelUserUpdateMsgID
	sqlChannelUserUnread
	sqlChannelUserChannels
	sqlChannelUserMembers
	sqlChannelUserlDisable
)

var (
	sqlChannelUser = []string {
		`CREATE TABLE IF NOT EXISTS chat.channel_user (
			channel_id	INTEGER UNSIGNED NOT NULL DEFAULT 0,
			user_id		BIGINT UNSIGNED NOT NULL DEFAULT 0,
			role 		TINYINT UNSIGNED NOT NULL DEFAULT 0,
			last_msg_id	BIGINT UNSIGNED NOT NULL DEFAULT 0,
            status      TINYINT UNSIGNED NOT NULL DEFAULT 0,
			INDEX (channel_id, user_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO chat.channel_user(channel_id,user_id,role) VALUES (?,?,?)`,
		`UPDATE chat.channel_user SET last_msg_id = ? WHERE user_id = ? AND channel_id = ? LIMIT 1`,
		`SELECT last_msg_id FROM chat.channel_user WHERE channel_id = ? AND user_id = ? LOCK IN SHARE MODE`,
		`SELECT channel_id FROM chat.channel_user WHERE user_id = ? LOCK IN SHARE MODE`,
		`SELECT user_id, role FROM chat.channel_user WHERE channel_id = ? LOCK IN SHARE MODE`,
		`UPDATE channel_user SET status = 1 WHERE channel_id = ?`,
	}
)

// Create table channel_user.
func (cu *channelUserServiceProvider) Create() error {
	_, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserCreateTable])

	return err
}

// Insert a raw.
func (cu *channelUserServiceProvider) Insert(chanID uint32, userID string, role uint8) (uint64, error) {
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return 0, err
	}

	result, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserInsert], chanID, userId, role)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvalidInsert
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

// UpdateMsgID update the lastMessageID by userID.
func (cu *channelUserServiceProvider) UpdateMsgID(msgID, userID string, chanID uint32) error {
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return err
	}

	result, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserUpdateMsgID], msgID, userId, chanID)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidInsert
	}

	return nil
}

// UnreadMsgID returns the lastMessageID by userID.
func (cu *channelUserServiceProvider) UnreadMsgID(channelID uint32, userID string) (uint64, error) {
	var (
		lastMsgID uint64
	)

	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return 0, err
	}

	if err := cu.store.db.QueryRow(sqlChannelUser[sqlChannelUserUnread], channelID, userId).Scan(&lastMsgID); err != nil {
		return 0, err
	}

	return lastMsgID, nil
}

// ChannelsByUserID returns the channelIDs by userID.
func (cu *channelUserServiceProvider) ChannelsByUserID(userID string) (*[]uint32, error) {
	var (
		single uint32
		result []uint32
	)

	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, err
	}

	rows, err := cu.store.db.Query(sqlChannelUser[sqlChannelUserChannels], userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&single); err != nil {
			return nil, err
		}

		result = append(result, single)
	}

	return &result, nil
}

// MemberByChanID return UserIDs by ChannelID.
func (cu *channelUserServiceProvider) MemberByChanID(channelID uint32) ([]*channelUser, error) {
	var (
		userID uint64
		role uint8
		result []*channelUser
	)

	rows, err := cu.store.db.Query(sqlChannelUser[sqlChannelUserMembers], channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&userID, &role); err != nil {
			return nil, err
		}

		single := &channelUser{
			UserID: userID,
			Role: role,
		}

		result = append(result, single)
	}

	return result, nil
}

func (cu *channelUserServiceProvider) Remove(channelID uint32) error {
	result, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserlDisable], channelID)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidDisable
	}

	return nil
}
