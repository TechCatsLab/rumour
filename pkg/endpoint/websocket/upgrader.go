/*
 * Revision History:
 *     Initial: 2018/07/24        Tong Yuehong
 */

package websocket

import (
	"fmt"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/rumour/identify"
	"github.com/TechCatsLab/rumour/pkg/conn"
	"github.com/TechCatsLab/rumour/pkg/log"
)

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

	c := conn.NewConn(ep.logic.Hub, ws, id)
	ep.logic.Hub.ConnectionManager.Add(c)

	c.Start()
	return nil
}
