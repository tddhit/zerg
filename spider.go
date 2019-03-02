package zerg

import (
	"runtime"
	"time"

	"github.com/tddhit/tools/log"
)

type Parser interface {
	Name() string
	Parse(rsp *Response) ([]*Item, []*Request)
}

type spider struct {
	parsers        map[string]Parser
	reqToEngineC   chan<- *Request
	itemToEngineC  chan<- *Item
	rspFromEngineC <-chan *Response
	poolC          chan struct{}
	seeds          chan *Request
}

func newSpider(
	reqToEngineC chan<- *Request,
	itemToEngineC chan<- *Item,
	rspFromEngineC <-chan *Response,
) *spider {

	return &spider{
		parsers:        make(map[string]Parser),
		reqToEngineC:   reqToEngineC,
		itemToEngineC:  itemToEngineC,
		rspFromEngineC: rspFromEngineC,
		poolC:          make(chan struct{}, runtime.GOMAXPROCS(runtime.NumCPU())),
		seeds:          make(chan *Request, 1000),
	}
}

func (s *spider) addParser(p Parser) {
	if _, ok := s.parsers[p.Name()]; !ok {
		s.parsers[p.Name()] = p
	} else {
		log.Warnf("parser[%s] is already exist!", p.Name())
	}
}

func (s *spider) addSeed(url, parser string) error {
	req, err := NewRequest(url, parser)
	if err != nil {
		log.Error(err)
		return err
	}
	s.seeds <- req
	return nil
}

func (s *spider) addRequest(req *Request) {
	s.seeds <- req
}

func (s *spider) start() {
	go func() {
		for req := range s.seeds {
			if req == nil {
				continue
			}
			if _, ok := s.parsers[req.Parser]; ok {
				s.reqToEngineC <- req
			} else {
				log.Errorf("req[url:%s, parser:%s] has no corresponding parser",
					req.RawURL, req.Parser)
			}
		}
	}()
	for {
		rsp := <-s.rspFromEngineC
		if rsp == nil {
			continue
		}
		if parser, ok := s.parsers[rsp.Parser]; !ok {
			log.Errorf("rsp[url:%s, parser:%s] has no corresponding parser",
				rsp.RawURL, rsp.Parser)
		} else {
			s.poolC <- struct{}{}
			go func(rsp *Response) {
				start := time.Now()
				items, reqs := parser.Parse(rsp)
				rsp.Body.Close()
				end := time.Now()
				elapsed := end.Sub(start)
				log.Debugf("parse %s spend %dms, req:%d, item:%d\n",
					rsp.RawURL, elapsed/1000000, len(reqs), len(items))
				for _, item := range items {
					if item == nil {
						continue
					}
					item.RawURL = rsp.RawURL
					s.itemToEngineC <- item
				}
				for _, req := range reqs {
					if req == nil {
						continue
					}
					if _, ok := s.parsers[req.Parser]; ok {
						s.reqToEngineC <- req
					}
				}
				<-s.poolC
			}(rsp)
		}
	}
}
