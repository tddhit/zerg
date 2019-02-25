package crawler

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/types"
)

type Proxy struct {
	Addr   string
	Valid  bool
	Added  bool
	Client *http.Client
	CloseC chan struct{}
}

type DefaultCrawler struct {
	workersPerProxy int
	reqC            chan *types.Request
	rspC            chan *types.Response
	pool            map[string]*types.Proxy
}

func NewDefaultCrawler(workersPerProxy int) *DefaultCrawler {
	return &DefaultCrawler{
		workersPerProxy: workersPerProxy,
		reqC:            make(chan *types.Request, 10000),
		pool:            make(map[string]*types.Proxy),
	}
}

func (c *DefaultCrawler) Name() string {
	return "DEFAULT_Crawler"
}

func (c *DefaultCrawler) Dispatch(req *types.Request) {
	select {
	case c.reqC <- req:
		return
	default:
		c.rspC <- &types.Response{
			RawURL: req.RawURL,
			err:    fmt.Errorf("crawler[%s] dispatch failed.", c.Name()),
		}
	}
}

func (c *DefaultCrawler) CrawlLoop(rspC chan<- *types.Response) {
	c.rspC = rspC
	for i := 0; i < c.workersPerProxy; i++ {
		go c.crawl(nil)
	}
}

func (c *DefaultCrawler) crawl(closeC <-chan struct{}) {
	for {
		select {
		case req := <-c.reqC:
			start := time.Now()
			rsp, err := client.Do(req.Request)
			if err != nil {
				log.Errorf("Failed Crawl %s %s\n!", req.RawURL, err)
				time.Sleep(2 * time.Second)
				continue
			}
			end := time.Now()
			elapsed := end.Sub(start)
			if rsp != nil {
				log.Infof("crawl %s(%s) spend %dms\n",
					req.RawURL, rsp.Response.Status, elapsed/1000000)
				irsp := &types.Response{
					Response: rsp,
					RawURL:   req.RawURL,
					Parser:   req.Parser,
				}
				d.rspC <- irsp
			}
		case <-closeC:
			return
		}
	}
}

func (c *DefaultCrawler) AddProxy(addr string) error {
	proxy := &types.Proxy{
		Addr:   addr,
		Valid:  true,
		Added:  true,
		CloseC: make(chan struct{}),
	}
	p, err := url.Parse(proxy.Addr)
	if err != nil {
		log.Error(err)
		proxy.Valid = false
		return err
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		DisableKeepAlives:  true,
		Dial: (&net.Dialer{
			Timeout: 2000 * time.Millisecond,
		}).Dial,
	}
	tr.Proxy = http.ProxyURL(p)
	proxy.Client = &http.Client{
		Transport: tr,
		Timeout:   15000 * time.Millisecond,
	}
	for i := 0; i < c.workersPerProxy; i++ {
		go c.crawl(proxy.CloseC)
	}
}

func (c *DefaultCrawler) RemoveProxy(addr string) error {
	if p, ok := c.pool[addr]; ok {
		close(p.CloseC)
		return nil
	} else {
		return ErrNotFound
	}
}
