package queuer

type DefaultQueuer struct {
	queue chan *Request
}

func NewDefaultQueuer() *DefaultQueuer {
	q := &DefaultQueuer{
		queue: make(chan *Request, 1000),
	}
}

func (q *DefaultQueuer) Push(req *Request) {
	q.queue <- req
}

func (q *DefaultQueuer) Pop() *Request {
	return <-q.queue
}
