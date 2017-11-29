package pipeline

import (
	"github.com/tddhit/zerg/pipeline/writer"
	"github.com/tddhit/zerg/types"
	"github.com/tddhit/zerg/util"
)

type Writer interface {
	Name() string
	Write(item *types.Item)
}

type Pipeline struct {
	writers            map[string]Writer
	itemFromEngineChan <-chan *types.Item
}

func NewPipeline(itemFromEngineChan <-chan *types.Item) *Pipeline {
	p := &Pipeline{
		writers:            make(map[string]Writer),
		itemFromEngineChan: itemFromEngineChan,
	}
	p.writers["DEFAULT_WRITER"] = writer.NewConsoleWriter("DEFAULT_WRITER")
	return p
}

func (p *Pipeline) AddWriter(writer Writer) *Pipeline {
	if _, ok := p.writers[writer.Name()]; !ok {
		p.writers[writer.Name()] = writer
	} else {
		util.LogFatal("writer[%s] is already exist!", writer.Name())
	}
	return p
}

func (p *Pipeline) Go() {
	go func() {
		for {
			item := <-p.itemFromEngineChan
			var writer Writer
			if w, ok := p.writers[item.Writer]; ok {
				writer = w
			} else {
				writer = p.writers["DEFAULT_WRITER"]
			}
			writer.Write(item)
		}
	}()
}
