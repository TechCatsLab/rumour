/*
 * Revision History:
 *     Initial: 2018/08/09        Tong Yuehong
 */

package api

import (
	"github.com/TechCatsLab/rumour/pkg/core"
	"github.com/TechCatsLab/apix/http/server"
)

type API struct {
	hub *core.Hub
}

func NewAPI(h *core.Hub) *API{
	return &API {
		hub: h,
	}
}

func Register(r *server.Router, h *core.Hub) {
	if r == nil {
		panic("API.Register: router is nil")
	}

	a := NewAPI(h)

	r.Post("/api/v1/single/unread", a.GetUnReadMessage)
	r.Post("/api/v1/single/history", a.GetRecord)

	r.Post("/api/v1/channel/joined", a.GetChannels)
	r.Post("/api/v1/channel/members", a.GetMembers)
	r.Post("/api/v1/channel/unread", a.GetChannelUnRead)
}
