package spider

import (
	"log"

	"github.com/tddhit/zerg/types"
)

type Parser interface {
	Parse(rsp *types.Response) (*types.Item, []*types.Request)
}

type Spider struct {
	Parser
	Name              string
	reqToEngineChan   chan<- *types.Request
	itemToEngineChan  chan<- *types.Item
	rspFromEngineChan <-chan *types.Response
	seeds             chan *types.Request
}

func NewSpider(name string, parser Parser) *Spider {
	s := &Spider{
		Parser: parser,
		Name:   name,
		seeds:  make(chan *types.Request, 100),
	}
	return s
}

func (s *Spider) AddSeed(url string) {
	req, _ := types.NewRequest(url, s.Name)
	select {
	case s.seeds <- req:
	default:
		log.Println("Warning: chan is full, discard!")
	}
}

/*
func (s *Spider) AssociateWriter(writer Writer) {
	s.associatedWriter = writer
}
*/

func (s *Spider) SetupChan(reqToEngineChan chan<- *types.Request,
	itemToEngineChan chan<- *types.Item, rspFromEngineChan <-chan *types.Response) {
	s.reqToEngineChan = reqToEngineChan
	s.itemToEngineChan = itemToEngineChan
	s.rspFromEngineChan = rspFromEngineChan
}

func (s *Spider) Go() {
	go func() {
		for req := range s.seeds {
			s.reqToEngineChan <- req
		}
	}()
	go func() {
		for {
			rsp := <-s.rspFromEngineChan
			if rsp != nil {
				item, reqs := s.Parse(rsp)
				if item != nil {
					s.itemToEngineChan <- item
				}
				for _, req := range reqs {
					s.reqToEngineChan <- req
				}
			}
		}
	}()
}
