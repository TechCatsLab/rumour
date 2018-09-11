/*
 * Revision History:
 *     Initial: 2018/08/24        Tong Yuehong
 */

package api

import (
	"errors"

	"github.com/TechCatsLab/apix/http/server"
	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/constants"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
	"github.com/TechCatsLab/rumour/response"
)

const (
	ChannelOwner = 2
)

var (
	errInvalidChannelOwner = errors.New("disable Channel: not the owner")
)

func (a *API) JoinChannel(c *server.Context) error {
	var (
		req struct {
			ChanID uint32 `json:"chan_id"`
			UserID string `json:"user_id"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	id, err := mysql.StoreService.Store().ChannelUser().Insert(req.ChanID, req.UserID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	err = a.hub.JoinChannel(req.UserID, req.ChanID)
	if err != nil {
		log.Error(err)
		return err
	}

	return response.WriteStatusAndIDJSON(c, constants.ErrSucceed, id)
}

func (a *API) ChangeRole(c *server.Context) error {
	var (
		req struct {
			OwnerID string `json:"owner_id"`
			ChanID  uint32 `json:"chan_id"`
			UserID  string `json:"user_id"`
			Role    uint8  `json:"role"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if req.OwnerID != "" {
		role, err := mysql.StoreService.Store().ChannelUser().GetRole(req.ChanID, req.OwnerID)
		if err != nil {
			log.Error(err)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}

		if role != ChannelOwner {
			log.Error(errInvalidChannelOwner)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}
	}

	err := mysql.StoreService.Store().ChannelUser().ChangeRole(req.ChanID, req.UserID, req.Role)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func (a *API) LeaveChannel(c *server.Context) error {
	var (
		req struct {
			ChanID  uint32 `json:"chan_id"`
			UserID  string `json:"user_id"`
			OwnerID string `json:"owner_id"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if req.OwnerID != "" {
		role, err := mysql.StoreService.Store().ChannelUser().GetRole(req.ChanID, req.UserID)
		if err != nil {
			log.Error(err)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}

		if role != ChannelOwner {
			log.Error(errInvalidChannelOwner)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}
	}

	err := mysql.StoreService.Store().ChannelUser().Exit(req.ChanID, req.UserID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}

func (a *API) UpdateMsgID(c *server.Context) error {
	var (
		req struct {
			ChanID uint32 `json:"chan_id"`
			MsgID  uint64 `json:"msg_id"`
			UserID string `json:"user_id"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	err := mysql.StoreService.Store().ChannelUser().UpdateMsgID(req.ChanID, req.MsgID, req.UserID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, nil)
}
