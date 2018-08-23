/*
 * Revision History:
 *     Initial: 2018/08/09        Tong Yuehong
 */

package api

import (
	"github.com/TechCatsLab/apix/http/server"

    log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/response"
	"github.com/TechCatsLab/rumour/constants"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
)

// Insert a channel.
func (a *API) Insert(c *server.Context) error {
	var (
		info struct {
			Name string `json:"name"`
			Title string `json:"title"`
		}
	)

	if err := c.JSONBody(&info); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	id, err := mysql.StoreService.Store().Channel().Insert(info.Name, info.Title)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndIDJSON(c, constants.ErrSucceed, id)
}

// Disable represents that a channel is deleted.
func (a *API) Disable(c *server.Context) error {
	var (
		info struct {
			ID uint32 `json:"id"`
		}
	)

	if err := c.JSONBody(&info); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err := mysql.StoreService.Store().Channel().Disable(info.ID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// GetMembers get members of a channel by channelID.
func (a *API) GetMembers(c *server.Context) error {
	var (
		info struct {
			ChanID uint32 `json:"chan_id"`
		}
	)

	if err := c.JSONBody(&info); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	members, err := mysql.StoreService.Store().ChannelUser().MemberByChanID(info.ChanID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, members)
}

// GetChannels get channels which the user joined by userID.
func (a *API) GetChannels(c *server.Context) error {
	var (
		info struct {
			UserID string `json:"user_id"`
		}
	)

	if err := c.JSONBody(&info); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	channels, err := mysql.StoreService.Store().ChannelUser().ChannelsByUserID(info.UserID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, channels)
}
