/*
 * Revision History:
 *     Initial: 2018/07/25        Tong Yuehong
 */

package store

type (
	channelServiceProvider struct {
		store *Store
	}

	Channel struct {
		Id    uint
		Name  string
		Title string
	}
)

// Create channels table.
func (c *channelServiceProvider) Create() error {
	_, err := c.store.DB.Exec(sqlStmts[SQLChannelCreateTable])

	return err
}

// Insert a Channel.
func (c *channelServiceProvider) Insert(name, title string) error {
	result, err := c.store.DB.Exec(sqlStmts[SQLChannelInsert], name, title)

	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidInsert
	}

	return nil
}

// Disable represent delete a Channel.
func (c *channelServiceProvider) Disable(channelID uint) error {
	result, err := c.store.DB.Exec(sqlStmts[SQLChannelDisable], channelID)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidDisable
	}

	return nil
}

// QueryByID query a Channel by ChannelID.
func (c *channelServiceProvider) QueryByID(channelID uint) (*Channel, error) {
	var channelInfo Channel

	if err := c.store.DB.QueryRow(sqlStmts[SQLChannelQueryByID], channelID).Scan(&channelInfo); err != nil {
		return nil, err
	}

	return &channelInfo, nil
}

// QueryByName returns a Channel named channelName.
func (c *channelServiceProvider) QueryByName(channelName string) (*Channel, error) {
	var channelInfo Channel

	if err := c.store.DB.QueryRow(sqlStmts[SQLChannelQueryByName], channelName).Scan(&channelInfo); err != nil {
		return nil, err
	}

	return &channelInfo, nil
}

// QueryExist query channels existed.
func (c *channelServiceProvider) QueryExist() ([]*Channel, error) {
	var (
		id uint
		name string
		title string
		channels []*Channel
	)

	rows, err := c.store.DB.Query(sqlStmts[SQLChannelQueryExist])
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
		id uint
		name string
		title string
		channels []*Channel
	)

	rows, err := c.store.DB.Query(sqlStmts[SQLChannelQueryDisabled])
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
