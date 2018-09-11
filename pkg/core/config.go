/*
 * Revision History:
 *     Initial: 2018/07/06        Tong Yuehong
 */

package core

// Config represents the configure when new a hub.
type Config struct {
	IncomingMessageQueueSize int

	DispatcherQueueSize int
	DispatcherWorkers   int
}

// NewConfig ensure the configure.
func NewConfig(fn ...func(*Config)) *Config {
	c := &Config{}

	for _, f := range fn {
		f(c)
	}

	return c
}

// Create a hub.
func (c *Config) Create() *Hub {
	return newHub(c)
}
