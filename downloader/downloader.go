package downloader

import (
	"sync"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/types"
)

type Crawler interface {
	Name() string
	Dispatch(req *types.Request)
	CrawlLoop(rspC chan<- *types.Response)
	AddProxy(addr string) error
	RemoveProxy(addr string) error
}

type Downloader struct {
	crawlers sync.Map
	reqC     <-chan *types.Request
	rspC     chan<- *types.Response
}

func NewDownloader(
	reqC <-chan *types.Request,
	rspC chan<- *types.Response,
) *Downloader {

	d := &Downloader{
		reqC: reqC,
		rspC: rspC,
	}
	return d
}

func (d *Downloader) RegisterCrawler(c Crawler) {
	if _, loaded := d.crawlers.LoadOrStore(c.Name(), c); loaded {
		log.Warnf("crawler[%s] is already exist!", c.Name())
		return
	}
	go c.CrawlLoop(d.rspC)
}

func (d *Downloader) AddProxy(addr string) *Downloader {
	d.crawlers.Range(func(key, value interface{}) bool {
		c := value.(Crawler)
		c.AddProxy(addr)
		return true
	})
	return d
}

func (d *Downloader) RemoveProxy(addr string) *Downloader {
	d.crawlers.Range(func(key, value interface{}) bool {
		c := value.(Crawler)
		c.RemoveProxy(addr)
		return true
	})
	return d
}

func (d *Downloader) Go() {
	go func() {
		for {
			req := <-d.reqC
			if c, ok := d.crawlers.Load(req.Crawler); ok {
				c.(Crawler).Dispatch(req)
			} else {
				c, _ = d.crawlers.Load("DEFAULT_Crawler")
				c.(Crawler).Dispatch(req)
			}
		}
	}()
}
