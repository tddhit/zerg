package spider

import "container/list"

type Parser interface {
	Parse(rsp *Response) (*Item, []*Request)
}

type Spider struct {
	Parser
	Name              string
	reqToEngineChan   chan<- *Request
	itemToEngineChan  chan<- *Item
	rspFromEngineChan <-chan *Response
	reqs              chan *Request
}

func NewSpider(name string, parser Parser) *Spider {
	s := &Spider{
		Parser: parser,
		Name:   name,
		reqs:   make(chan *Request, 1000),
	}
}

func (s *Spider) AddSeed(url string) {
	req, _ := http.NewRequest("GET", url, nil)
	ireq := &engine.Request{
		Request: req,
		RawURL:  url,
	}
	select {
	case s.reqs <- ireq:
	default:
		log.Println("Warning: chan is full, discard!")
	}
}

/*
func (s *Spider) AssociateWriter(writer Writer) {
	s.associatedWriter = writer
}
*/

func (s *Spider) SetupChan(reqToEngineChan chan<- *Request,
	itemToEngineChan chan<- *Item, rspFromEngineChan <-chan *Response) {
	s.reqToEngineChan = reqToEngineChan
	s.itemToEngineChan = itemToEngineChan
	s.rspFromEngineChan = rspFromEngineChan
}

func (s *Spider) Go() {
	go func() {
		for {
			req := <-s.reqs
			if req != nil {
				reqToEngineChan <- req
			}
		}
	}()
	go func() {
		for {
			rsp := <-rspFromEngineChan
			if rsp != nil {
				item, reqs := s.Parse(rsp)
				if item != nil {
					itemToEngineChan <- item
				}
				for _, req := range reqs {
					reqToEngineChan <- req
				}
			}
		}
	}()
}
