/**
 * Revision History:
 *     Initial: 2018/05/30        Tong Yuehong
 */

package rumour

// Connection -
type Connection interface {
	Start()
	Stop()
	Identify() (string, error)
	Send(*Message) error
}

type Queue interface {
	Put(*Message) error
	Get() (*Message, error)
	Close()
}

// Authenticator -
type Authenticator interface {
	Authenticate(string) (string, error)
}
