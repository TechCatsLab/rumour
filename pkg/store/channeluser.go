/*
 * Revision History:
 *     Initial: 2018/07/26        Tong Yuehong
 */

package store

type (
	channelUserServiceProvider struct {
		store *Store
	}

	channelUser struct {
		Id            uint64
		ChannelID     uint
		UserID        uint64
		Role          uint8
		LastMessageID uint64
	}
)

// Create table channel_user.
func (cu *channelUserServiceProvider) Create() error {
	_, err := cu.store.DB.Exec(sqlStmts[SQLChannelUserCreateTable])

	return err
}

// Insert a raw.
func (cu *channelUserServiceProvider) Insert(chanID uint, userID uint64, role uint8) error {
	result, err := cu.store.DB.Exec(sqlStmts[SQLChannelUserInsert])
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidInsert
	}

	return nil
}

// UpdateMsgID update the lastMessageID by userID.
func (cu *channelUserServiceProvider) UpdateMsgID(msgID, userID uint64, chanID uint) error {
	result, err := cu.store.DB.Exec(sqlStmts[SQLChannelUserUpdateMsgID])
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidInsert
	}

	return nil
}

// UnreadMsgID returns the lastMessageID by userID.
func (cu *channelUserServiceProvider) UnreadMsgID(channelID uint, userID uint64) (*uint64, error) {
	var (
		lastMsgID uint64
	)
	if err := cu.store.DB.QueryRow(sqlStmts[SQLChannelUserUnread], channelID, userID).Scan(&lastMsgID); err != nil {
		return nil, err
	}

	return &lastMsgID, nil
}

// ChannelsByUserID returns the channelIDs by userID.
func (cu *channelUserServiceProvider) ChannelsByUserID(userID uint64) (*[]uint, error) {
	var (
		single uint
		result []uint
	)
	rows, err := cu.store.DB.Query(sqlStmts[SQLChannelUserChannels])
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
func (cu *channelUserServiceProvider) MemberByChanID(channelID uint) ([]*channelUser, error) {
	var (
		userID uint64
		role uint8
		result []*channelUser
	)

	rows, err := cu.store.DB.Query(sqlStmts[SQLChannelUserMembers])
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
