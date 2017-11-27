package pipeline

type Writer interface {
	Write(item *Item)
}

type Pipeline struct {
	Writer
	itemFromEngineChan <-chan *Item
}

func NewPipeline(itemFromEngineChan <-chan *Item) *Pipeline {
	p := &Pipeline{
		Writer:             writer.NewWriter(),
		itemFromEngineChan: itemFromEngineChan,
	}
}

func (p *Pipeline) Go() {
	go func() {
		for {
			item := <-p.itemFromEngineChan
			p.Write(item)
		}
	}()
}
