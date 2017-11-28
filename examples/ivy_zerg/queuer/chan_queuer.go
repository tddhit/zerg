package queuer

import (
	"github.com/tddhit/zerg/types"
)

type ChanQueuer struct {
	queue chan *types.Request
}

func NewChanQueuer() *ChanQueuer {
	q := &ChanQueuer{
		queue: make(chan *types.Request, 1000),
	}
	return q
}

func (q *ChanQueuer) Push(req *types.Request) {
	q.queue <- req
}

func (q *ChanQueuer) Pop() *types.Request {
	return <-q.queue
}
