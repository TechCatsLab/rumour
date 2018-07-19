/*
 * Revision History:
 *     Initial: 2018/06/18        Tong Yuehong
 */

package websocket

import (
	"fmt"
	"net/http"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/gorilla/websocket"

	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/identify"
	"github.com/TechCatsLab/rumour/pkg/conn"
	"github.com/TechCatsLab/rumour/pkg/hub"
	"github.com/TechCatsLab/rumour/pkg/log"
	"github.com/TechCatsLab/rumour/pkg/parser"
)

const (
	socketBufferSize = 1024
)

type Endpoint struct {
	hub      rumour.Hub
	server   *server.Entrypoint
	router   *server.Router
	upgrader *websocket.Upgrader
}

var (
	fakeID        = 100
	configuration = &server.Configuration{
		Address: ":8088",
	}
)

func NewEndpoint() rumour.Endpoint {
	hubConfig := hub.New(hub.HubQueueSize, hub.DispatcherScheduler)
	ep := &Endpoint{
		hub: hubConfig.Create(),
		server: server.NewEntrypoint(configuration, nil),
		router: server.NewRouter(),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  socketBufferSize,
			WriteBufferSize: socketBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	ep.router.Get("/ws", ep.websocketHandler)
	return ep
}

func (ep *Endpoint) websocketHandler(ctx *server.Context) error {
	resp := ctx.Response()
	req := ctx.Request()

	fakeID = fakeID + 1
	id := identify.Identify(fmt.Sprintf("%d", fakeID))
	fmt.Println("id:", fakeID)


	ws, err := ep.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Error("[Websocket Endpoint] Upgrade err:", log.Err(err))
		resp.Write([]byte("Upgrade err"))
		return err
	}

	c := conn.NewConn(ep.hub, ws, parser.NewPacketParser(), id)
	ep.hub.ConnManager().Add(c)

	c.Start()
	return nil
}

func (ep *Endpoint) Serve() {
	ep.server.Start(ep.router.Handler())
	ep.server.Run()
}
