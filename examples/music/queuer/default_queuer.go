package queuer

import (
	"time"

	"github.com/tddhit/zerg/types"
)

type DefaultQueuer struct {
	queue chan *types.Request
}

func NewDefaultQueuer() *DefaultQueuer {
	q := &DefaultQueuer{
		queue: make(chan *types.Request, 1000000),
	}
	return q
}

func (q *DefaultQueuer) Push(req *types.Request) {
	q.queue <- req
}

func (q *DefaultQueuer) Pop() *types.Request {
	time.Sleep(20 * time.Millisecond)
	return <-q.queue
}
