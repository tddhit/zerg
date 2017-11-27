package scheduler

import (
	"github.com/tddhit/zerg/scheduler/queuer"
)

type Queuer interface {
	Push(req *Request)
	Pop() *Request
}

type Scheduler struct {
	Queuer
	reqToEngineChan   chan<- *Request
	reqFromEngineChan <-chan *Request
}

func NewScheduler(reqFromEngineChan <-chan *Reuqest, reqToEngineChan chan<- *Request) *Scheduler {
	s := &Scheduler{
		Queuer:            queuer.NewDefaultQueuer(),
		reqToEngineChan:   reqToEngineChan,
		reqFromEngineChan: reqFromEngineChan,
	}
}

func (s *Scheduler) SetQueuer(q Queuer) {
	s.Queuer = q
}

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
