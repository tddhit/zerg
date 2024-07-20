package crawler

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg"
)

type proxy struct {
	addr   string
	client *http.Client
	closeC chan struct{}
}

type options struct {
	concurrency int
	interval    time.Duration
	dialTimeout time.Duration
	readTimeout time.Duration
}

type Option func(*options)

var defaultOptions = options{
	concurrency: 4,
	interval:    2 * time.Second,
	dialTimeout: 5 * time.Second,
	readTimeout: 15 * time.Second,
}

func WithConcurrency(n int) Option {
	return func(o *options) {
		o.concurrency = n
	}
}

func WithDialTimeout(t time.Duration) Option {
	return func(o *options) {
		o.dialTimeout = t
	}
}

func WithReadTimeout(t time.Duration) Option {
	return func(o *options) {
		o.readTimeout = t
	}
}

func WithInterval(t time.Duration) Option {
	return func(o *options) {
		o.interval = t
	}
}

type DefaultCrawler struct {
	opt          options
	reqC         chan *zerg.Request
	rspToEngineC chan<- *zerg.Response
	closeC       chan struct{}
	poolMutex    sync.Mutex
	pool         map[string]*proxy
}

func NewDefaultCrawler(opts ...Option) *DefaultCrawler {
	opt := defaultOptions
	for _, o := range opts {
		o(&opt)
	}
	return &DefaultCrawler{
		opt:    opt,
		reqC:   make(chan *zerg.Request, 1000),
		closeC: make(chan struct{}),
		pool:   make(map[string]*proxy),
	}
}

func (c *DefaultCrawler) Name() string {
	return "DEFAULT_CRAWLER"
}

func (c *DefaultCrawler) Dispatch(req *zerg.Request) {
	select {
	case c.reqC <- req:
		return
	default:
		c.rspToEngineC <- &zerg.Response{
			Request: req,
			Err: fmt.Errorf("crawler[url:%s, crawler:%s] dispatch failed.",
				req.RawURL, c.Name()),
		}
	}
}

func (c *DefaultCrawler) CrawlLoop(rspToEngineC chan<- *zerg.Response) {
	c.rspToEngineC = rspToEngineC
	for i := 0; i < c.opt.concurrency; i++ {
		go c.crawl(c.newClient(nil), nil)
	}
}

func (c *DefaultCrawler) crawl(client *http.Client, proxyCloseC <-chan struct{}) {
	for {
		select {
		case req := <-c.reqC:
			start := time.Now()
			rsp, err := client.Do(req.Request)
			if err != nil {
				log.Errorf("Failed Crawl %s %s\n!", req, err)
				time.Sleep(c.opt.interval)
				c.rspToEngineC <- &zerg.Response{
					Request: req,
					Err:     fmt.Errorf("%s crawl failed: %s.", c.Name(), err),
				}
				continue
			}
			end := time.Now()
			elapsed := end.Sub(start)
			log.Infof("crawl %s(%s) spend %dms\n",
				req.RawURL, rsp.Status, elapsed/1000000)
			irsp := &zerg.Response{
				Request:  req,
				Response: rsp,
			}
			c.rspToEngineC <- irsp
			time.Sleep(c.opt.interval)
		case <-c.closeC:
			close(c.reqC)
			c.reqC = nil
			return
		case <-proxyCloseC:
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
				Timeout: c.opt.dialTimeout,
			}).Dial,
			Proxy: proxyURL,
		},
		Timeout: c.opt.readTimeout,
	}
}

func (c *DefaultCrawler) AddProxy(addr string) {
	p, err := url.Parse(addr)
	if err != nil {
		log.Error(err)
		return
	}

	c.poolMutex.Lock()
	defer c.poolMutex.Unlock()

	if _, ok := c.pool[addr]; ok {
		log.Errorf("proxy[%s] already exist.", addr)
	} else {
		conn, err := net.DialTimeout(
			"tcp",
			strings.TrimPrefix(addr, "http://"),
			2*time.Second,
		)
		if err != nil {
			log.Errorf("dial proxy[%s] failed. err:%s", addr, err)
			return
		}
		conn.Close()
		c.pool[addr] = &proxy{
			addr:   addr,
			closeC: make(chan struct{}),
			client: c.newClient(http.ProxyURL(p)),
		}
		for i := 0; i < c.opt.concurrency; i++ {
			go c.crawl(c.pool[addr].client, c.pool[addr].closeC)
		}
	}
}

func (c *DefaultCrawler) RemoveProxy(addr string) {
	c.poolMutex.Lock()
	defer c.poolMutex.Unlock()

	if p, ok := c.pool[addr]; ok {
		delete(c.pool, addr)
		close(p.closeC)
	} else {
		log.Errorf("proxy[%s] not found.", addr)
	}
}

func (c *DefaultCrawler) Close() {
	close(c.closeC)
}
