package downloader

import (
	"time"

	"github.com/tddhit/tools/log"

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
		Crawler:           &crawler.HTTPCrawler{},
		reqFromEngineChan: reqFromEngineChan,
		rspToEngineChan:   rspToEngineChan,
	}
	return d
}

func (d *Downloader) Go() {
	go func() {
		for {
			req := <-d.reqFromEngineChan
			go func(req *types.Request) {
				start := time.Now()
				rsp := d.Crawl(req)
				end := time.Now()
				elapsed := end.Sub(start)
				if rsp != nil {
					log.Debugf("crawl %s(%s) spend %dms\n", req.RawURL, rsp.Response.Status, elapsed/1000000)
					if rsp.Response.StatusCode == 200 {
						rsp.RawURL = req.RawURL
						rsp.Parser = req.Parser
						d.rspToEngineChan <- rsp
					} else {
						rsp.Body.Close()
					}
				}
			}(req)
		}
	}()
}
