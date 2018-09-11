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

// Create a channel.
func (a *API) CreateChannel(c *server.Context) error {
	var (
		req struct {
			Name string `json:"name"`
			Title string `json:"title"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	id, err := mysql.StoreService.Store().Channel().Insert(req.Name, req.Title)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndIDJSON(c, constants.ErrSucceed, id)
}

// Disable represents that a channel is deleted.
func (a *API) DisableChannel(c *server.Context) error {
	var (
		req struct {
			UserID string `json:"user_id"`
			ChanID uint32 `json:"chan_id"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err := mysql.StoreService.Store().Channel().Disable(req.ChanID, req.UserID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

// ListMembers get members of a channel by channelID.
func (a *API) ListMembers(c *server.Context) error {
	var (
		req struct {
			ChanID uint32 `json:"chan_id"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	members, err := mysql.StoreService.Store().ChannelUser().MemberByChanID(req.ChanID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, members)
}

// ListChannels get channels which the user joined by userID.
func (a *API) ListChannels(c *server.Context) error {
	var (
		req struct {
			UserID string `json:"user_id"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	channels, err := mysql.StoreService.Store().ChannelUser().ChannelsByUserID(req.UserID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, channels)
}
