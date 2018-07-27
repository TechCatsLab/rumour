/*
 * Revision History:
 *     Initial: 2018/07/19        Tong Yuehong
 */

package server

import (
	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/pkg/core"
	"github.com/TechCatsLab/rumour/pkg/endpoint/websocket"
	"github.com/TechCatsLab/rumour/pkg/hub"
)

func CreateWebsocketServer() rumour.Endpoint {
	return websocket.NewEndpoint(
		core.NewLogic(*hub.NewConfig(hub.HubQueueSize, hub.DispatcherScheduler).Create()),
	)
}
