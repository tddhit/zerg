package scheduler

import (
	"github.com/tddhit/zerg/scheduler/queuer"
	"github.com/tddhit/zerg/types"
)

type Queuer interface {
	Push(req *types.Request)
	Pop() *types.Request
}

type Scheduler struct {
	Queuer
	reqToEngineChan   chan<- *types.Request
	reqFromEngineChan <-chan *types.Request
}

func NewScheduler(reqFromEngineChan <-chan *types.Request, reqToEngineChan chan<- *types.Request) *Scheduler {
	s := &Scheduler{
		Queuer:            queuer.NewDefaultQueuer(),
		reqToEngineChan:   reqToEngineChan,
		reqFromEngineChan: reqFromEngineChan,
	}
	return s
}

/*
func (s *Scheduler) SetQueuer(q Queuer) {
	s.Queuer = q
}
*/

func (s *Scheduler) Go() {
	go func() {
		for {
			req := s.Pop()
			if req != nil {
				s.reqToEngineChan <- req
			}
		}
	}()
	go func() {
		for {
			req := <-s.reqFromEngineChan
			if req != nil {
				s.Push(req)
			}
		}
	}()
}
