package downloader

import (
	"github.com/tddhit/zerg/downloader/crawler"
	"github.com/tddhit/zerg/types"
)

type Crawler interface {
	Crawl(req *types.Request) *types.Response
}

type Downloader struct {
	Crawler
	reqFromEngineChan <-chan *types.Request
	rspToEngineChan   chan<- *types.Response
}

func NewDownloader(reqFromEngineChan <-chan *types.Request, rspToEngineChan chan<- *types.Response) *Downloader {
	d := &Downloader{
		Crawler:           crawler.NewHTTPCrawler(),
		reqFromEngineChan: reqFromEngineChan,
		rspToEngineChan:   rspToEngineChan,
	}
	return d
}

func (d *Downloader) Go() {
	go func() {
		for {
			req := <-d.reqFromEngineChan
			rsp := d.Crawl(req)
			if rsp != nil && rsp.Response.StatusCode == 200 {
				rsp.RawURL = req.RawURL
				rsp.Parser = req.Parser
				d.rspToEngineChan <- rsp
			}
		}
	}()
}
