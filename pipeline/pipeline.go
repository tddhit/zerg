package pipeline

import (
	"github.com/tddhit/zerg/pipeline/writer"
	"github.com/tddhit/zerg/types"
	"github.com/tddhit/zerg/util"
)

type Writer interface {
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
	p.writers["DEFAULT_WRITER"] = writer.NewConsoleWriter()
	return p
}

func (p *Pipeline) AssociateWriter(spiderName string, writer Writer) {
	if _, ok := p.writers[spiderName]; !ok {
		p.writers[spiderName] = writer
	} else {
		util.LogWarn("spider[%s] already has a writer!", spiderName)
	}
}

func (p *Pipeline) Go() {
	go func() {
		for {
			item := <-p.itemFromEngineChan
			var writer Writer
			if w, ok := p.writers[item.Spider]; ok {
				writer = w
			} else {
				writer = p.writers["DEFAULT_WRITER"]
			}
			writer.Write(item)
		}
	}()
}
