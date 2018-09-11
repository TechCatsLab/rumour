/*
 * Revision History:
 *     Initial: 2018/07/26        Tong Yuehong
 */

package mysql

import (
	"strconv"
	"errors"
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
	sqlChannelUserDisable
	sqlChannelUserGetRole
	sqlChannelUserChangeRole
	sqlChannelUserExit

	channelOwner = 2
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
		`INSERT INTO chat.channel_user(channel_id,user_id) VALUES (?,?)`,
		`UPDATE chat.channel_user SET last_msg_id = ? WHERE user_id = ? AND channel_id = ? LIMIT 1`,
		`SELECT last_msg_id FROM chat.channel_user WHERE channel_id = ? AND user_id = ? LOCK IN SHARE MODE`,
		`SELECT channel_id FROM chat.channel_user WHERE user_id = ? AND status = 0 LOCK IN SHARE MODE`,
		`SELECT user_id, role FROM chat.channel_user WHERE channel_id = ? LOCK IN SHARE MODE`,
		`UPDATE channel_user SET status = 1 WHERE channel_id = ?`,
		`SELECT role FROM chat.channel_user WHERE channel_id = ? AND user_id = ? AND status = 0`,
		`UPDATE channel_user SET role = ? WHERE channel_id = ? AND user_id = ? AND status = 0`,
		`UPDATE channel_user SET status = 1 WHERE channel_id = ? AND user_id = ?`,
	}
)

var (
	errInvalidCUInsert = errors.New("insert channel_user: insert affected 0 rows")
	errInvalidCUDisable = errors.New("disable channel_user: disable affected 0 rows")
	errUpdateLastMsgID = errors.New("updateMsgID channel_user: update message_id affected 0 rows")
	errUpdateRole = errors.New("updateRole channel_user: update role affected 0 rows")
)

// Create table channel_user.
func (cu *channelUserServiceProvider) Create() error {
	_, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserCreateTable])

	return err
}

// Insert a raw.
func (cu *channelUserServiceProvider) Insert(channelID uint32, userID string) (uint64, error) {
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return 0, err
	}

	result, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserInsert], channelID, userId)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvalidCUInsert
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

// UpdateMsgID update the lastMessageID by userID.
func (cu *channelUserServiceProvider) UpdateMsgID(channelID uint32, msgID uint64, userID string) error {
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return err
	}

	result, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserUpdateMsgID], msgID, userId, channelID)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errUpdateLastMsgID
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
	result, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserDisable], channelID)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidCUDisable
	}

	return nil
}

func (cu *channelUserServiceProvider) GetRole(chanID uint32, userID string) (uint8, error) {
	var (
		role uint8
	)

	if err := cu.store.db.QueryRow(sqlChannelUser[sqlChannelUserGetRole], chanID, userID).Scan(&role); err != nil {
		return 44, err
	}

	return role, nil
}

func (cu *channelUserServiceProvider) ChangeRole(chanID uint32, userID string, role uint8) error {
	result, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserChangeRole], role, chanID, userID)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errUpdateRole
	}

	return nil
}

func (cu *channelUserServiceProvider) Exit(channelID uint32, userID string) error {
	result, err := cu.store.db.Exec(sqlChannelUser[sqlChannelUserExit], channelID, userID)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errUpdateRole
	}

	return nil
}
