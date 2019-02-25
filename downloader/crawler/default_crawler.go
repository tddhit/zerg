package crawler

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
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
	rspC            chan<- *types.Response
	poolMutex       sync.Mutex
	pool            map[string]*Proxy
}

func NewDefaultCrawler(workersPerProxy int) *DefaultCrawler {
	return &DefaultCrawler{
		workersPerProxy: workersPerProxy,
		reqC:            make(chan *types.Request, 10000),
		pool:            make(map[string]*Proxy),
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
			Request: req,
			Err:     fmt.Errorf("crawler[%s] dispatch failed.", c.Name()),
		}
	}
}

func (c *DefaultCrawler) CrawlLoop(rspC chan<- *types.Response) {
	c.rspC = rspC
	for i := 0; i < c.workersPerProxy; i++ {
		go c.crawl(c.newClient(nil), nil)
	}
}

func (c *DefaultCrawler) crawl(client *http.Client, closeC <-chan struct{}) {
	for {
		select {
		case req := <-c.reqC:
			start := time.Now()
			rsp, err := client.Do(req.Request)
			if err != nil {
				log.Errorf("Failed Crawl %s %s\n!", req.RawURL, err)
				time.Sleep(2 * time.Second)
				c.rspC <- &types.Response{
					Request: req,
					Err:     fmt.Errorf("crawl failed: %s.", c.Name(), err),
				}
				continue
			}
			end := time.Now()
			elapsed := end.Sub(start)
			log.Infof("crawl %s(%s) spend %dms\n",
				req.RawURL, rsp.Status, elapsed/1000000)
			irsp := &types.Response{
				Response: rsp,
				Request:  req,
			}
			c.rspC <- irsp
			time.Sleep(time.Second)
		case <-closeC:
			return
		}
	}
}

func (c *DefaultCrawler) newClient(
	proxyURL func(*http.Request) (*url.URL, error),
) *http.Client {

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
			DisableKeepAlives:  true,
			Dial: (&net.Dialer{
				Timeout: 5000 * time.Millisecond,
			}).Dial,
			Proxy: proxyURL,
		},
		Timeout: 15000 * time.Millisecond,
	}
}

func (c *DefaultCrawler) AddProxy(addr string) error {
	p, err := url.Parse(addr)
	if err != nil {
		log.Error(err)
		return err
	}
	proxy := &Proxy{
		Addr:   addr,
		Valid:  true,
		Added:  true,
		CloseC: make(chan struct{}),
		Client: c.newClient(http.ProxyURL(p)),
	}

	c.poolMutex.Lock()
	if _, ok := c.pool[addr]; ok {
		c.poolMutex.Unlock()
		return fmt.Errorf("proxy[%s] already exist.", addr)
	} else {
		c.pool[addr] = proxy
	}
	c.poolMutex.Unlock()

	for i := 0; i < c.workersPerProxy; i++ {
		go c.crawl(proxy.Client, proxy.CloseC)
	}
	return nil
}

func (c *DefaultCrawler) RemoveProxy(addr string) error {
	c.poolMutex.Lock()
	defer c.poolMutex.Unlock()

	if p, ok := c.pool[addr]; ok {
		delete(c.pool, addr)
		close(p.CloseC)
		return nil
	} else {
		return fmt.Errorf("proxy[%s] not found.", addr)
	}
}
