/*
 * Revision History:
 *     Initial: 2018/07/25        Tong Yuehong
 */

package mysql

import (
	"time"
	"errors"
	"database/sql"
	"context"
)

type (
	channelServiceProvider struct {
		store *store
	}

	Channel struct {
		Id        uint32
		Name      string
		Title     string
		CreatedAt time.Time
		Status    uint8
	}
)

const (
	sqlChannelCreateTable                     = iota
	sqlChannelInsert
	sqlChannelDisable
	sqlChannelQueryByID
	sqlChannelQueryByName
	sqlChannelQueryExist
	sqlChannelQueryDisabled
)

var (
	sqlChannel = []string{
		`CREATE TABLE IF NOT EXISTS chat.channels (
			id 			INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
			name 		VARCHAR(255) UNIQUE NOT NULL DEFAULT ' ',
			title 		VARCHAR(512) NOT NULL DEFAULT ' ',
			created_at 	DATETIME NOT NULL DEFAULT current_timestamp,
			status      TINYINT UNSIGNED NOT NULL DEFAULT 0,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO channels(name,title) VALUES (?,?)`,
		`UPDATE channels SET status = 1 WHERE Id = ? LIMIT 1`,
		`SELECT id,name,title FROM chat.channels WHERE id = ? AND status = 0 LOCK IN SHARE MODE`,
		`SELECT id,name,title FROM chat.channels WHERE name = ? AND status = 0 LOCK IN SHARE MODE`,
		`SELECT id,name,title FROM chat.channels WHERE status = 0 LOCK IN SHARE MODE`,
		`SELECT id,name,title FROM chat.channels WHERE status = 1 LOCK IN SHARE MODE`,
	}
)

var (
	errInvalidChannelInsert = errors.New("insert Channel: insert affected 0 rows")
	errInvalidChannelDisable = errors.New("disable Channel: disable affected 0 rows")
	errInvalidChannelOwner = errors.New("disable Channel: not the owner")
)

// Create channels table.
func (c *channelServiceProvider) Create() error {
	_, err := c.store.db.Exec(sqlChannel[sqlChannelCreateTable])

	return err
}

// Insert a Channel.
func (c *channelServiceProvider) Insert(name, title string) (uint32, error) {
	result, err := c.store.db.Exec(sqlChannel[sqlChannelInsert], name, title)

	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvalidChannelInsert
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

// Disable represent delete a Channel.
func (c *channelServiceProvider) Disable(channelID uint32, userID string) (err error) {
	role, err := c.store.channelUser.GetRole(channelID, userID)
	if role != channelOwner {
		return errInvalidChannelOwner
	}

	tx, err := c.store.db.BeginTx(context.Background(), &sql.TxOptions{})
	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	result, err := tx.Exec(sqlChannel[sqlChannelDisable], channelID)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChannelDisable
	}

	remove, err := tx.Exec(sqlChannelUser[sqlChannelUserDisable], channelID)
	if err != nil {
		return err
	}

	if affected, _ := remove.RowsAffected(); affected == 0 {
		return errInvalidCUDisable
	}

	return nil
}

// QueryByID query a Channel by ChannelID.
func (c *channelServiceProvider) QueryByID(channelID uint32) (*Channel, error) {
	var (
		id          uint32
		name        string
		title       string
		channelInfo Channel
	)

	if err := c.store.db.QueryRow(sqlChannel[sqlChannelQueryByID], channelID).Scan(&id, &name, &title); err != nil {
		return nil, err
	}

	channelInfo = Channel{
		Id:    id,
		Name:  name,
		Title: title,
	}

	return &channelInfo, nil
}

// QueryByName returns a Channel named channelName.
func (c *channelServiceProvider) QueryByName(channelName string) (*Channel, error) {
	var channelInfo Channel

	if err := c.store.db.QueryRow(sqlChannel[sqlChannelQueryByName], channelName).Scan(&channelInfo); err != nil {
		return nil, err
	}

	return &channelInfo, nil
}

// QueryExist query channels existed.
func (c *channelServiceProvider) QueryExist() ([]*Channel, error) {
	var (
		id       uint32
		name     string
		title    string
		channels []*Channel
	)

	rows, err := c.store.db.Query(sqlChannel[sqlChannelQueryExist])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &title); err != nil {
			return nil, err
		}

		channel := &Channel{
			Id:    id,
			Name:  name,
			Title: title,
		}

		channels = append(channels, channel)
	}

	return channels, nil
}

// QueryDisabled query channels disabled.
func (c *channelServiceProvider) QueryDisabled() ([]*Channel, error) {
	var (
		id       uint32
		name     string
		title    string
		channels []*Channel
	)

	rows, err := c.store.db.Query(sqlChannel[sqlChannelQueryDisabled])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, name, title); err != nil {
			return nil, err
		}

		channel := &Channel{
			Id:    id,
			Name:  name,
			Title: title,
		}

		channels = append(channels, channel)
	}

	return channels, nil
}
