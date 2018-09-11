/*
 * Revision History:
 *     Initial: 2018/09/03        Tong Yuehong
 */

package api

import (
	"github.com/TechCatsLab/apix/http/server"
	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/response"
	"github.com/TechCatsLab/rumour/constants"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
)

func (a *API) ListChannelUnRead(c *server.Context) error {
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

	if req.MsgID == 0 {
		lastMsgID,err := mysql.StoreService.Store().ChannelUser().UnreadMsgID(req.ChanID, req.UserID)
		if err != nil {
			log.Error(err)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}

		req.MsgID = lastMsgID
	}

	messages, err := mysql.StoreService.Store().ChannelMessage().Unread(req.ChanID, req.MsgID)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, messages)
}

