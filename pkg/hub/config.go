/*
 * Revision History:
 *     Initial: 2018/07/06        Tong Yuehong
 */

package hub

type Config struct {
	IncomingMessageQueueSize int

	DispatcherQueueSize int
	DispatcherWorkers   int
}

func NewConfig(fn ...func(*Config)) *Config {
	c := &Config{}

	for _, f := range fn {
		f(c)
	}

	return c
}

func (c *Config) Create() *Hub {
	return newHub(c)
}
