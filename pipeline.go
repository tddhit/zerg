package zerg

import (
	"github.com/tddhit/tools/log"
)

type Writer interface {
	Name() string
	Write(item *Item)
}

type pipeline struct {
	writers map[string]Writer
	itemC   <-chan *Item
}

func newPipeline(itemC <-chan *Item) *pipeline {
	return &pipeline{
		writers: make(map[string]Writer),
		itemC:   itemC,
	}
}

func (p *pipeline) addWriter(w Writer) {
	if _, ok := p.writers[w.Name()]; !ok {
		p.writers[w.Name()] = w
	} else {
		log.Warnf("writer[%s] is already exist!", w.Name())
	}
	log.Debug("!!!:", w.Name())
}

func (p *pipeline) start() {
	for {
		item := <-p.itemC
		for _, name := range item.Writers {
			if w, ok := p.writers[name]; ok {
				w.Write(item)
			} else {
				log.Errorf("item[url:%s, writer:%s] has no corresponding writer",
					item.RawURL, name)
			}
		}
	}
}
