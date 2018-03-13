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

type HTTPCrawler struct {
}

func (c *HTTPCrawler) Crawl(req *types.Request) *types.Response {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		DisableKeepAlives:  true,
		Dial: (&net.Dialer{
			Timeout: 2000 * time.Millisecond,
		}).Dial,
	}
	if req.Proxy != "" {
		proxy, err := url.Parse(req.Proxy)
		if err != nil {
			return nil
		}
		tr.Proxy = http.ProxyURL(proxy)
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   5000 * time.Millisecond,
	}
	rsp, err := client.Do(req.Request)
	if err != nil {
		log.Errorf("Failed Crawl %s %s\n!", req.RawURL, err)
		return nil
	}
	irsp := &types.Response{
		Response: rsp,
	}
	return irsp
}
