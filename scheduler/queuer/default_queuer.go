package queuer

import (
	"github.com/tddhit/zerg/types"
)

type DefaultQueuer struct {
	queue chan *types.Request
}

func NewDefaultQueuer() *DefaultQueuer {
	q := &DefaultQueuer{
		queue: make(chan *types.Request, 1000),
	}
	return q
}

func (q *DefaultQueuer) Push(req *types.Request) {
	q.queue <- req
}

func (q *DefaultQueuer) Pop() *types.Request {
	return <-q.queue
}
