package zerg

type Queuer interface {
	Push(req *Request)
	Pop() *Request
}

type scheduler struct {
	queuer         Queuer
	reqToEngineC   chan<- *Request
	reqFromEngineC <-chan *Request
}

func newScheduler(
	reqFromEngineC <-chan *Request,
	reqToEngineC chan<- *Request,
) *scheduler {

	return &scheduler{
		reqToEngineC:   reqToEngineC,
		reqFromEngineC: reqFromEngineC,
	}
}

func (s *scheduler) setQueuer(q Queuer) {
	s.queuer = q
}

func (s *scheduler) start() {
	go func() {
		for {
			req := s.queuer.Pop()
			if req != nil {
				s.reqToEngineC <- req
			}
		}
	}()
	for {
		req := <-s.reqFromEngineC
		if req != nil {
			s.queuer.Push(req)
		}
	}
}
