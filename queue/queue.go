package queue

import (
	"sync"
)

type queueStorage struct {
	mu    sync.Mutex
	queue []string
}

func (q *queueStorage) Enqueue(value string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue = append(q.queue, value)
}

func (q *queueStorage) Dequeue() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) < 1 {
		return "", false
	}
	data := q.queue[0]
	q.queue[0] = ""
	q.queue = q.queue[1:]
	return data, true
}

func (q *queueStorage) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.queue) == 0
}

func New(values []string) *queueStorage {
	// copy urls to queue
	q := &queueStorage{}
	q.queue = make([]string, len(values))
	copy(q.queue, values)

	return q
}
