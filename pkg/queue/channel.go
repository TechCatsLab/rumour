/*
 * Revision History:
 *     Initial: 2018/06/27        Tong Yuehong
 */

package queue

import (
	"sync"

	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/pkg/log"
)

type queue struct {
	queue chan rumour.Message
	stop sync.Once
}

// NewQueue new a queue.
func NewChannelQueue(size int) rumour.Queue {
	return &queue {
		queue: make(chan rumour.Message, size),
	}
}

// Put a message on queue.
func (q *queue) Put(message rumour.Message) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error("[Queue Put] Put message err", log.Err(err.(error)))
			// metric
		}
	}()

	q.queue <- message

	return nil
}

// Get a message from queue.
func (q *queue) Get() (rumour.Message, error) {
	return <-q.queue, nil
}

// Close the queue for receiving messages.
func (q *queue) Close() {
	q.stop.Do(func() {
		close(q.queue)
	})
}
