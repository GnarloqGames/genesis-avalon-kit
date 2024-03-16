package couchbase

import (
	"log/slog"
	"sync"
)

var queue *Queue

type Queue struct {
	mx *sync.Mutex

	queue    []any
	Incoming chan any
}

func GetQueue() (*Queue, error) {
	if queue == nil {
		queue = &Queue{
			mx:       &sync.Mutex{},
			queue:    make([]any, 0),
			Incoming: make(chan any),
		}

		db, err := Get()
		if err != nil {
			return nil, err
		}

		go func() {
			for {
				select {
				case item := <-queue.Incoming:
					queue.mx.Lock()
					queue.queue = append(queue.queue, item)
					queue.mx.Unlock()
				default:
					if len(queue.queue) == 0 {
						continue
					}

					queue.mx.Lock()
					item := queue.queue[0]
					queue.queue = queue.queue[1:]

					if err := db.upsert(item); err != nil {
						slog.Error("failed to upsert item", "error", err)

						if err != ErrInvalidItemType {
							queue.queue = append(queue.queue, item)
						}
					}
					queue.mx.Unlock()
				}
			}
		}()
	}

	return queue, nil
}
