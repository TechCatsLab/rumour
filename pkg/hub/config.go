/*
 * Revision History:
 *     Initial: 2018/07/06        Tong Yuehong
 */

package hub

import (
	"github.com/TechCatsLab/rumour"
)

type Config struct {
	IncomingMessageQueueSize int

	DispatcherQueueSize int
	DispatcherWorkers   int
}

func New(fn ...func(*Config)) *Config {
	c := &Config{}

	for _, f := range fn {
		f(c)
	}

	return c
}

func (c *Config) Create() rumour.Hub {
	return newHub(c)
}
