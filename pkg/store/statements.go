/*
 * Revision History:
 *     Initial: 2018/07/25        Tong Yuehong
 */

package store

import (
	"errors"
)

const (
	SQLChannelCreateTable = iota
	SQLChannelInsert
	SQLChannelDisable
	SQLChannelQueryByID
	SQLChannelQueryByName
	SQLChannelQueryExist
	SQLChannelQueryDisabled

	SQLSingleMessageCreateTable
	SQLSingleMessageInsert
	SQLSingleMessageQueryRecord
	SQLSingleMessageUnreadByUserID
	SQLSingleMessageUnreadByMsgID

	SQLChannelMessageCreateTable
	SQLChannelMessageInsert
	SQLChannelMessageUnread

	SQLChannelUserCreateTable
	SQLChannelUserInsert
	SQLChannelUserUpdateMsgID
	SQLChannelUserUnread
	SQLChannelUserChannels
	SQLChannelUserMembers
)

var (
	errInvalidInsert = errors.New("store Channel: insert affected 0 rows")
	errInvalidDisable = errors.New("store Channel: disable affected 0 rows")

	sqlStmts = []string{
		`CREATE TABLE IF NOT EXISTS channels (
			id 			integer unsigned NOT NULL AUTO_INCREMENT,
			name 		varchar(255) NOT NULL DEFAULT ' ',
			title 		varchar(512) NOT NULL DEFAULT ' ',
			created_at 	datetime NOT NULL DEFAULT current_timestamp,
			status      tinyint unsigned NOT NULL DEFAULT 0,
			INDEX name,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8_bin;`,
		`INSERT INTO channels(name,title) VALUES (?,?)`,
		`UPDATE channels SET status = 1 WHERE Id = ? LIMIT 1`,
		`SELECT id,name,title FROM channels WHERE id = ? AND status = 0 LOCK IN SHARE MODE`,
		`SELECT id,name,title FROM channels WHERE name = ? AND status = 0 LOCK IN SHARE MODE`,
		`SELECT id,name,title FROM channels WHERE status = 0 LOCK IN SHARE MODE`,
		`SELECT id,name,title FROM channels WHERE status = 1 LOCK IN SHARE MODE`,

		`CREATE TABLE IF NOT EXISTS single_message (
			Id 			bigint unsigned NOT NULL,
			source_id 	bigint unsigned NOT NULL DEFAULT 0,
			target_id 	bigint unsigned NOT NULL DEFAULT 0,
			kind 		integer unsigned NOT NULL DEFAULT 0,
			content     text NOT NULL DEFAULT '',
			created_at 	datetime NOT NULL DEFAULT current_timestamp,
			arrived 	bool NOT NULL DEFAULT false,
			INDEX (source_id, target_id),
			PRIMARY KEY (Id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8_bin;`,
		`INSERT INTO single_message(Id, source_id,target_id,kind,Content) VALUES (?,?,?,?,?)`,
		`SELECT id,source_id,target_id,kind,content,created_at FROM single_message WHERE source_id = ? AND target_id = ? OR source_id = ? AND target_id = ? LOCK IN SHARE MODE`,
		`SELECT id,source_id,target_id,kind,content,created_at FROM single_message WHERE source_id = ? AND target_id = ? AND arrived = false LOCK IN SHARE MODE`,
		`SELECT id,source_id,target_id,kind,content,created_at FROM single_message WHERE id > ? AND target_id = ? AND source_id = ? AND arrived = false LOCK IN SHARE MODE`,

		`CREATE TABLE IF NOT EXISTS channel_message (
			id 			bigint unsigned NOT NULL,
			channel_id	integer unsigned NOT NULL DEFAULT 0,
			source_id 	bigint unsigned NOT NULL DEFAULT 0,
			kind 		integer unsigned NOT NULL DEFAULT 0,
			content		text NOT NULL DEFAULT '',
			created_at 	datetime NOT NULL DEFAULT current_timestamp,
			INDEX (id, channel_id),
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8_bin;`,
		`INSERT INTO channel_message(Id,source_id,channel_id,kind,Content) VALUES (?,?,?,?,?)`,
		`SELECT id,source_id,kind,content,created_at FROM channel_message WHERE channel_id = ? AND id > ?`,

		`CREATE TABLE IF NOT EXISTS channel_user (
			channel_id	integer unsigned NOT NULL DEFAULT 0,
			user_id		bigint unsigned NOT NULL DEFAULT 0,
			role 		tinyint unsigned NOT NULL DEFAULT 0,
			last_msg_id	bigint unsigned NOT NULL DEFAULT 0,
			INDEX (channel_id, user_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8_bin;`,
		`INSERT INTO channel_user(channel_id,user_id,role) VALUES (?,?,?)`,
		`UPDATE channel_user SET last_msg_id = ? WHERE user_id = ? AND channel_id = ? LIMIT 1`,
		`SELECT last_msg_id FROM channel_user WHERE channel_id = ? AND user_id = ? LOCK IN SHARE MODE`,
		`SELECT channel_id FROM channel_user WHERE user_id = ? LOCK IN SHARE MODE`,
		`SELECT user_id, role FROM channel_user WHERE channel_id = ? LOCK IN SHARE MODE`,
	}
)
