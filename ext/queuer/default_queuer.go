package queuer

import (
	"sync"

	"github.com/tddhit/zerg"
)

type DefaultQueuer struct {
	sync.RWMutex
	visited map[string]struct{}
	queue   chan *zerg.Request
}

func NewDefaultQueuer() *DefaultQueuer {
	return &DefaultQueuer{
		visited: make(map[string]struct{}),
		queue:   make(chan *zerg.Request, 10000),
	}
}

func (q *DefaultQueuer) Push(req *zerg.Request) {
	q.RLock()
	if _, ok := q.visited[req.RawURL]; ok {
		q.RUnlock()
		return
	}
	q.RUnlock()

	q.Lock()
	if _, ok := q.visited[req.RawURL]; ok {
		q.Unlock()
		return
	}
	q.visited[req.RawURL] = struct{}{}
	q.Unlock()

	q.queue <- req
}

func (q *DefaultQueuer) Pop() *zerg.Request {
	return <-q.queue
}
