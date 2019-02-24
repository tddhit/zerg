package downloader

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/tddhit/tools/log"

	"github.com/tddhit/zerg/downloader/crawler"
	"github.com/tddhit/zerg/types"
)

type Crawler interface {
	Name() string
	Crawl(req *types.Request) *types.Response
}

type Downloader struct {
	crawlers map[string]Crawler

	//代理
	poolMutex  sync.RWMutex
	proxyPool  []string
	proxyIndex int32

	reqFromEngineChan <-chan *types.Request
	rspToEngineChan   chan<- *types.Response
}

func NewDownloader(reqFromEngineChan <-chan *types.Request, rspToEngineChan chan<- *types.Response) *Downloader {
	d := &Downloader{
		crawlers:          make(map[string]Crawler),
		reqFromEngineChan: reqFromEngineChan,
		rspToEngineChan:   rspToEngineChan,
	}
	d.crawlers["DEFAULT_HTTPCrawler"] = &crawler.HTTPCrawler{}
	return d
}

func (d *Downloader) AddProxy(proxy string) *Downloader {
	d.poolMutex.Lock()
	d.proxyPool = append(d.proxyPool, proxy)
	d.poolMutex.Unlock()
	return d
}

func (d *Downloader) AddCrawler(c Crawler) *Downloader {
	if _, ok := d.crawlers[c.Name()]; !ok {
		d.crawlers[c.Name()] = c
	} else {
		log.Warnf("crawler[%s] is already exist!", c.Name())
	}
	return d
}

func (d *Downloader) Go() {
	go func() {
		for {
			req := <-d.reqFromEngineChan
			go func(req *types.Request) {
				var cr Crawler
				if c, ok := d.crawlers[req.Crawler]; ok {
					cr = c
				} else {
					cr = d.crawlers["DEFAULT_HTTPCrawler"]
				}
				start := time.Now()
				if req.Proxy == "" && len(d.proxyPool) > 0 {
					d.poolMutex.RLock()
					idx := atomic.AddInt32(&d.proxyIndex, 1)
					idx %= int32(len(d.proxyPool))
					req.Proxy = d.proxyPool[idx]
					d.poolMutex.RUnlock()
				}
				rsp := cr.Crawl(req)
				end := time.Now()
				elapsed := end.Sub(start)
				if rsp != nil {
					log.Infof("crawl %s(%s) spend %dms\n",
						req.RawURL, rsp.Response.Status, elapsed/1000000)
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
