/*
 * Revision History:
 *     Initial: 2018/05/22        Tong Yuehong
 */

package main

import "github.com/TechCatsLab/rumour/pkg/endpoint/websocket"

func main() {
	websocket.NewEndpoint().Serve()
}
