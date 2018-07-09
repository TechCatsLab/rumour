/**
 * Revision History:
 *     Initial: 2018/05/30        Tong Yuehong
 */

package rumour

// Hub -
type Hub interface {
	HubDispatcher() Dispatcher
	ConnManager() ConnectionManager
	Dispatch(message Message) error
}

// ConnectionManager -
type ConnectionManager interface {
	Add(Connection) error
	Remove(Connection) error
	Query(Identify) ([]Connection, error)
}

// Connection -
type Connection interface {
	Hub() Hub
	Start()
	Stop()
	Identify() Identify
	Send(Message) error
}

// Dispatcher -
type Dispatcher interface {
	Hub() Hub
	Dispatch(message Message) error
}

// Queue -
type Queue interface {
	Put(Message) error
	Get() (Message, error)
	Close()
}

// Parser -
type Parser interface {
	Parse([]byte) (Message, error)
}

// Message -
type Message interface {
	Source() Identify
	Target() Identify
	MessageType() interface{}
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

// Identify -
type Identify interface {
	Equal(Identify) bool
	Id() (string, error)
}

// Authenticator -
type Authenticator interface {
	Authenticate(string) (string, error)
}
