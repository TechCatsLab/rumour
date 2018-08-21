/*
 * Revision History:
 *     Initial: 2018/07/29        Tong Yuehong
 */

package websocket

import (
	"errors"
	"sync"

	"github.com/TechCatsLab/rumour"
)

var (
	errMessageTypeHandlerExists = errors.New("webSocketHandler: message type exists")
)

type webSocketHandler interface {
	Serve(*rumour.Message) error
}

type WebSocketRouter struct {
	handlers map[rumour.MessageType]webSocketHandler
	mu       sync.RWMutex
}

type HandleFunc struct {
	handlerFunc func(*rumour.Message) error
}

func (h *HandleFunc) Serve(m *rumour.Message) error {
	if err := h.handlerFunc(m); err != nil {
		return err
	}

	return nil
}

func (wr *WebSocketRouter) HandlerFunc(f func(*rumour.Message) error) webSocketHandler {
	return &HandleFunc{f}
}

func New() *WebSocketRouter {
	return &WebSocketRouter{
		handlers: make(map[rumour.MessageType]webSocketHandler),
	}
}

func (wr *WebSocketRouter) Handle(mtype rumour.MessageType, handler webSocketHandler) error {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if _, exists := wr.handlers[mtype]; !exists {
		wr.handlers[mtype] = handler
		return nil
	}

	return errMessageTypeHandlerExists
}

func (wr *WebSocketRouter) ServeWebSocket(message *rumour.Message) {
	if h, ok := wr.handlers[rumour.MessageType(message.Type)]; !ok {
		return
	} else {
		h.Serve(message)
	}
}
