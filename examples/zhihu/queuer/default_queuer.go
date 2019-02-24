package queuer

import (
	"time"

	"github.com/tddhit/zerg/types"
)

type DelayQueuer struct {
	queue chan *types.Request
}

func NewDelayQueuer() *DelayQueuer {
	q := &DelayQueuer{
		queue: make(chan *types.Request, 1000000),
	}
	return q
}

func (q *DelayQueuer) Push(req *types.Request) {
	q.queue <- req
}

func (q *DelayQueuer) Pop() *types.Request {
	time.Sleep(20 * time.Millisecond)
	return <-q.queue
}
