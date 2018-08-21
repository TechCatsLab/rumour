/*
 * Revision History:
 *     Initial: 2018/05/27        Tong Yuehong
 */

package conn

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/TechCatsLab/rumour"
	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/pkg/queue"
	"github.com/TechCatsLab/rumour/pkg/core"
)

var (
	ErrDifferentConn = errors.New("different connection")
	ErrInvalidID     = errors.New("invalid identify")
	ErrMessageSeq    = errors.New("seq error")
)

type conn struct {
	hub      *core.Hub
	seq      int
	ws       *websocket.Conn
	queue    rumour.Queue
	identify string
	shutdown chan struct{}
	stop     sync.Once
}

//NewConn - create a new Conn.
func NewConn(hub *core.Hub, ws *websocket.Conn, identify string) rumour.Connection {
	conn := &conn{
		ws:       ws,
		hub:      hub,
		shutdown: make(chan struct{}),
		queue:    queue.NewChannelQueue(1024),
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
			if err == ErrMessageSeq {
				continue
			}
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
		log.Error(err)
		c.Stop()
		return err
	}

	if c.identify != (msg.From) {
		log.Error("[conn handleRead] different connection")
		return ErrDifferentConn
	}

	id := c.hub.Generator().Get()
	msg.Content["id"] = id

	err = c.hub.Put(msg)
	if err != nil {
		return nil
	}

	ack := &rumour.Message{
		Seq:     msg.Seq,
		Type:    0,
		Content: make(map[string]interface{}),
	}

	ack.Content["id"] = id
	return c.Send(ack)
}

func (c *conn) handleWrite() error {
	msg, err := c.queue.Get()
	if err != nil {
		log.Error(err)
		return err
	}

	return c.ws.WriteJSON(msg)
}

// Receive message which is sent by client.
func (c *conn) receive() (*rumour.Message, error) {
	_, b, err := c.ws.ReadMessage()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	msg, err := c.parse(b)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if msg.Seq <= c.seq {
		return nil, ErrMessageSeq
	}

	c.seq = msg.Seq

	return msg, nil
}

func (c *conn) parse(data []byte) (*rumour.Message, error) {
	var m rumour.Message

	err := m.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	if m.From == "" || m.To == "" || m.Type == 0 {
		return nil, errors.New("param err")
	}

	return &m, nil
}

// Identify return the id.
func (c *conn) Identify() (string, error) {
	if c.identify == "" {
		return "", ErrInvalidID
	}
	return c.identify, nil
}

// Send a message to the queue saved message which is sent to client.
func (c *conn) Send(message *rumour.Message) error {
	return c.queue.Put(message)
}

// Stop the connection.
func (c *conn) Stop() {
	c.stop.Do(func() {
		c.hub.ConnectionManager.Remove(c)
		c.ws.Close()
		c.queue.Close()
		close(c.shutdown)
	})
}
