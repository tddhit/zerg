package zerg

import (
	"github.com/tddhit/tools/log"
)

type Crawler interface {
	Name() string
	Dispatch(req *Request)
	CrawlLoop(rspToEngineC chan<- *Response)
	AddProxy(addr string)
	RemoveProxy(addr string)
	Close()
}

type downloader struct {
	crawlers       map[string]Crawler
	reqFromEngineC <-chan *Request
	rspToEngineC   chan<- *Response
}

func newDownloader(
	reqFromEngineC <-chan *Request,
	rspToEngineC chan<- *Response,
) *downloader {

	return &downloader{
		crawlers:       make(map[string]Crawler),
		reqFromEngineC: reqFromEngineC,
		rspToEngineC:   rspToEngineC,
	}
}

func (d *downloader) addCrawler(c Crawler) {
	if _, ok := d.crawlers[c.Name()]; ok {
		log.Warnf("crawler[%s] is already exist!", c.Name())
		return
	} else {
		d.crawlers[c.Name()] = c
		go c.CrawlLoop(d.rspToEngineC)
	}
}

func (d *downloader) addProxy(addr string) {
	for _, c := range d.crawlers {
		c.AddProxy(addr)
	}
}

func (d *downloader) removeProxy(addr string) {
	for _, c := range d.crawlers {
		c.RemoveProxy(addr)
	}
}

func (d *downloader) start() {
	for {
		req := <-d.reqFromEngineC
		if c, ok := d.crawlers[req.Crawler]; ok {
			c.Dispatch(req)
		} else {
			log.Errorf("req[url:%s, crawler:%s] has no corresponding crawler",
				req.RawURL, req.Crawler)
		}
	}
}

func (d *downloader) close() {
	for _, c := range d.crawlers {
		c.Close()
	}
}
