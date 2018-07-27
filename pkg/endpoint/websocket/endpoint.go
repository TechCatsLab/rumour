/*
 * Revision History:
 *     Initial: 2018/06/18        Tong Yuehong
 */

package websocket

import (
	"net/http"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/gorilla/websocket"

	"github.com/TechCatsLab/rumour/pkg/core"
	"github.com/TechCatsLab/rumour"
)

const (
	socketBufferSize = 1024
)

type Endpoint struct {
	logic    *core.Logic
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

func NewEndpoint(logic *core.Logic) rumour.Endpoint {
	ep := &Endpoint{
		logic: logic,
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

func (ep *Endpoint) Serve() error {
	if err := ep.server.Start(ep.router.Handler()); err != nil {
		return err
	}

	ep.server.Run()
	return nil
}
