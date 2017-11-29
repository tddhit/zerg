package spider

import (
	"github.com/tddhit/zerg/types"
	"github.com/tddhit/zerg/util"
)

type Parser interface {
	Name() string
	Parse(rsp *types.Response) (*types.Item, []*types.Request)
}

type Spider struct {
	parsers           map[string]Parser
	reqToEngineChan   chan<- *types.Request
	itemToEngineChan  chan<- *types.Item
	rspFromEngineChan <-chan *types.Response
	seeds             chan *types.Request
}

func NewSpider(reqToEngineChan chan<- *types.Request,
	itemToEngineChan chan<- *types.Item, rspFromEngineChan <-chan *types.Response) *Spider {
	s := &Spider{
		parsers:           make(map[string]Parser),
		reqToEngineChan:   reqToEngineChan,
		itemToEngineChan:  itemToEngineChan,
		rspFromEngineChan: rspFromEngineChan,
		seeds:             make(chan *types.Request, 100),
	}
	return s
}

func (s *Spider) AddParser(parser Parser) *Spider {
	if _, ok := s.parsers[parser.Name()]; !ok {
		s.parsers[parser.Name()] = parser
	} else {
		util.LogWarn("parser[%s] is already exist!", parser.Name())
	}
	return s
}

func (s *Spider) AddSeed(url, parser string) *Spider {
	req, _ := types.NewRequest(url, parser)
	select {
	case s.seeds <- req:
	default:
		util.LogWarn("spider -> engine, chan is full, discard %s!", req.RawURL)
	}
	return s
}

func (s *Spider) Go() {
	go func() {
		for req := range s.seeds {
			if req != nil {
				if _, ok := s.parsers[req.Parser]; ok {
					s.reqToEngineChan <- req
				}
			}
		}
	}()
	go func() {
		for {
			rsp := <-s.rspFromEngineChan
			if rsp != nil {
				if parser, ok := s.parsers[rsp.Parser]; ok {
					item, reqs := parser.Parse(rsp)
					if item != nil {
						item.RawURL = rsp.RawURL
						util.LogDebug("%s\n", item.Writer)
						s.itemToEngineChan <- item
					}
					for _, req := range reqs {
						if req != nil {
							if _, ok := s.parsers[req.Parser]; ok {
								s.reqToEngineChan <- req
							}
						}
					}
				}
			}
		}
	}()
}
