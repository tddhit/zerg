package pipeline

import (
	"github.com/tddhit/zerg/pipeline/writer"
	"github.com/tddhit/zerg/types"
)

type Writer interface {
	Write(item *types.Item)
}

type Pipeline struct {
	Writer
	itemFromEngineChan <-chan *types.Item
}

func NewPipeline(itemFromEngineChan <-chan *types.Item) *Pipeline {
	p := &Pipeline{
		Writer:             writer.NewConsoleWriter(),
		itemFromEngineChan: itemFromEngineChan,
	}
	return p
}

func (p *Pipeline) Go() {
	go func() {
		for {
			item := <-p.itemFromEngineChan
			p.Write(item)
		}
	}()
}
