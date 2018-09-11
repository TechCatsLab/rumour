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

func (a *API) FetchUnreadMessages(c *server.Context) error {
	var (
		req struct {
			From  string `json:"from"`
			To    string `json:"to"`
			MsgID uint64 `json:"msg_id"`
		}

		err error

		messages interface{}
	)

	if err = c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	if req.MsgID == 0 {
		messages, err = mysql.StoreService.Store().SingleMessage().QueryUnreadByUserID(req.From, req.To)
		if err != nil {
			log.Error(err)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}
	} else {
		messages, err = mysql.StoreService.Store().SingleMessage().QueryUnreadByMsgID(req.From, req.To, req.MsgID)
		if err != nil {
			log.Error(err)
			return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
		}
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, messages)
}

func (a *API) ListHistoryMessages(c *server.Context) error {
	var (
		req struct {
			From string `json:"from"`
			To   string `json:"to"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrInvalidParam, nil)
	}

	messages, err := mysql.StoreService.Store().SingleMessage().QueryRecord(req.From, req.To)
	if err != nil {
		log.Error(err)
		return response.WriteStatusAndDataJSON(c, constants.ErrMysql, nil)
	}

	return response.WriteStatusAndDataJSON(c, constants.ErrSucceed, messages)
}
