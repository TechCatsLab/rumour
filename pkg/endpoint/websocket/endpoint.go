/*
 * Revision History:
 *     Initial: 2018/06/18        Tong Yuehong
 */

package websocket

import (
	"fmt"
	"net/http"
	_ "github.com/go-sql-driver/mysql"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/gorilla/websocket"

	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/pkg/conn"
	"github.com/TechCatsLab/rumour/pkg/core"
	"github.com/TechCatsLab/rumour/pkg/endpoint/api"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
)

const (
	socketBufferSize = 1024
)

type Endpoint struct {
	Hub      *core.Hub
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

func NewEndpoint(hub *core.Hub) rumour.Endpoint {
	ep := &Endpoint{
		Hub:    hub,
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

	api.Register(ep.router, ep.Hub)

	ep.Hub.ChannelManager.Load()

	ep.router.Get("/ws", ep.websocketHandler)
	return ep
}

func (ep *Endpoint) websocketHandler(ctx *server.Context) error {
	resp := ctx.Response()
	req := ctx.Request()

	fakeID = fakeID + 1 // TODO: fix
	id := fmt.Sprintf("%d", fakeID)

	ws, err := ep.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Error(err)
		resp.Write([]byte("Upgrade err"))
		return err
	}

	c := conn.NewConn(ep.Hub, ws, id)
	ep.Hub.ConnectionManager.Add(c)

	chans, err := mysql.StoreService.Store().ChannelUser().ChannelsByUserID(id)
	if err != nil {
		return err
	}

	if chans != nil {
		for _, channel := range *chans {
			single, err := mysql.StoreService.Store().Channel().QueryByID(channel)
			if err != nil {
				log.Error(err)
				return err
			}

			err = ep.Hub.ChannelManager.Add(single.Id, c)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}

	c.Start()
	return nil
}

func (ep *Endpoint) Serve() error {
	if err := ep.server.Start(ep.router.Handler()); err != nil {
		return err
	}

	ep.server.Run()
	return nil
}
