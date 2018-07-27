/**
 * Revision History:
 *     Initial: 2018/05/30        Tong Yuehong
 */

package rumour

// Connection -
type Connection interface {
	Start()
	Stop()
	Identify() Identify
	Send(Message) error
}

// Queue -
type Queue interface {
	Put(Message) error
	Get() (Message, error)
	Close()
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
