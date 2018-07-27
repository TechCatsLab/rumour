/*
 * Revision History:
 *     Initial: 2018/07/25        Tong Yuehong
 */

package core

import (
	"github.com/TechCatsLab/rumour/pkg/channel"
	"github.com/TechCatsLab/rumour/pkg/hub"
	"github.com/TechCatsLab/rumour/pkg/store"
)

type Logic struct {
	Hub            hub.Hub
	ChannelManager *channel.ChannelManager
	Store          *store.Store
}

func NewLogic(hub hub.Hub) *Logic {
	return &Logic{
		Hub: hub,
	}
}
