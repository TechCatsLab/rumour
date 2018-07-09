/*
 * Revision History:
 *     Initial: 2018/05/27        Tong Yuehong
 */

package conn

import (
	"sync"
	"errors"

	"github.com/gorilla/websocket"

	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/pkg/log"
	"github.com/TechCatsLab/rumour/pkg/queue"
	"github.com/TechCatsLab/rumour/constants"
)

var ErrDifferentConn = errors.New("different connection")
type conn struct {
	hub      rumour.Hub
	ws       *websocket.Conn
	parser   rumour.Parser
	queue    rumour.Queue
	identify rumour.Identify
	shutdown chan struct{}
	stop     sync.Once
}

//NewConn - create a new Conn.
func NewConn(hub rumour.Hub, ws *websocket.Conn, parser rumour.Parser, identify rumour.Identify) rumour.Connection {
	conn := &conn{
		ws:       ws,
		hub:      hub,
		parser:   parser,
		shutdown: make(chan struct{}),
		queue:    queue.NewChannelQueue(constants.ConnQueueSize),
		identify: identify,
	}

	return conn
}

func (c *conn) Start() {
	go c.readLoop()
	go c.writeLoop()
}

func (c *conn) readLoop() {
	for {
		if err := c.handleRead(); err != nil {
			c.Stop()
			break
		}

		select {
		case <-c.shutdown:
			return
		default:
		}
	}
}

func (c *conn) writeLoop() {
	for {
		select {
		case <-c.shutdown:
			return
		default:
		}

		if err := c.handleWrite(); err != nil {
			c.Stop()
			break
		}
	}
}

func (c *conn) handleRead() error {
	msg, err := c.receive()
	if err != nil {
		log.Error("[Connection handleRead] HandleRead err:", log.Err(err))
		c.Stop()
		return err
	}

	if !c.identify.Equal(msg.Source()) {
		log.Error("[conn handleRead] different connection")
		return ErrDifferentConn
	}

	return c.Hub().Dispatch(msg)
}

func (c *conn) handleWrite() error {
	msg, err := c.queue.Get()
	if err != nil {
		log.Error("[Connection handleWrite] HandleWrite err:", log.Err(err))
		return err
	}

	return c.ws.WriteJSON(msg)
}

// Hub return the hub.
func (c *conn) Hub() rumour.Hub {
	return c.hub
}

// Identify return the identify.
func (c *conn) Identify() rumour.Identify {
	return c.identify
}

// Receive message which is sent by client.
func (c *conn) receive() (rumour.Message, error) {
	_, b, err := c.ws.ReadMessage()
	if err != nil {
		log.Error("[Connection Receive]Can't read from websocket", log.Err(err))
		return nil, err
	}

	msg, err := c.parser.Parse(b)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Send a message to the client.
func (c *conn) Send(message rumour.Message) error {
	return c.queue.Put(message)
}

// Stop the connection.
func (c *conn) Stop() {
	c.stop.Do(func() {
		c.hub.ConnManager().Remove(c)
		c.ws.Close()
		c.queue.Close()
		close(c.shutdown)
	})
}
