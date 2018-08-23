/*
 * Revision History:
 *     Initial: 2018/08/09        Tong Yuehong
 */

package api

import (
	"github.com/TechCatsLab/apix/http/server"

	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/constants"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
	"github.com/TechCatsLab/rumour/response"
)

func (a *API) GetUnReadMessage(c *server.Context) error {
	var (
		info struct {
			From  string `json:"from"`
			To    string `json:"to"`
			MsgID uint64 `json:"msg_id"`
		}

		err error

		messages interface{}
	)

	if err = c.JSONBody(&info); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if info.MsgID == 0 {
		messages, err = mysql.StoreService.Store().SingleMessage().QueryUnreadByUserID(info.From, info.To)
		if err != nil {
			log.Error(err)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}
	} else {
		messages, err = mysql.StoreService.Store().SingleMessage().QueryUnreadByMsgID(info.From, info.To, info.MsgID)
		if err != nil {
			log.Error(err)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, messages)
}

func (a *API) GetRecord(c *server.Context) error {
	var (
		info struct {
			From string `json:"from"`
			To   string `json:"to"`
		}
	)

	if err := c.JSONBody(&info); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	messages, err := mysql.StoreService.Store().SingleMessage().QueryRecord(info.From, info.To)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, messages)
}

func (a *API) GetChannelUnRead(c *server.Context) error {
	var (
		info struct {
			ChanID uint32 `json:"chan_id"`
			MsgID  uint64 `json:"msg_id"`
			UserID string `json:"user_id"`
		}
	)

	if err := c.JSONBody(&info); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if info.MsgID == 0 {
		lastMsgID,err := mysql.StoreService.Store().ChannelUser().UnreadMsgID(info.ChanID, info.UserID)
		if err != nil {
			log.Error(err)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}

		info.MsgID = lastMsgID
	}

	messages, err := mysql.StoreService.Store().ChannelMessage().Unread(info.ChanID, info.MsgID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, messages)
}
