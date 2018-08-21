/*
 * Revision History:
 *     Initial: 2018/07/17        Tong Yuehong
 */

package core

import (
	"runtime"
)

func HubQueueSize(c *Config) {
	c.IncomingMessageQueueSize = 1024 * 16
}

func DispatcherScheduler(c *Config) {
	c.DispatcherWorkers = 2 * runtime.NumCPU()
	c.DispatcherQueueSize = 2 * c.DispatcherWorkers
}
