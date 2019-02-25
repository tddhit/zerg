package downloader

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/types"
)

type Crawler interface {
	Name() string
	Dispatch(req *types.Request)
	CrawlLoop(rspC chan<- *types.Response)
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
	d.crawlers.Range(func(key, value interface{}) {
		c := value.(Crawler)
		c.AddProxy(addr)
	})
}

func (d *Downloader) RemoveProxy(addr string) *Downloader {
	d.crawlers.Range(func(key, value interface{}) {
		c := value.(Crawler)
		c.RemoveProxy(addr)
	})
}

func (d *Downloader) Go() {
	go func() {
		for {
			req := <-d.reqC
			if c, ok := d.crawlers.Load(req.Crawler); ok {
				c.Dispatch(req)
			} else {
				c = d.crawlers.Load("DEFAULT_Crawler")
				c.Dispatch(req)
			}
		}
	}()
}
