package crawler

import (
	"crypto/tls"
	"net/http"

	"github.com/tddhit/tools/log"

	"github.com/tddhit/zerg/types"
)

type HTTPCrawler struct {
	*http.Client
}

func NewHTTPCrawler() *HTTPCrawler {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{
		Transport: tr,
	}
	crawler := &HTTPCrawler{
		Client: client,
	}
	return crawler
}

func (c *HTTPCrawler) Crawl(req *types.Request) *types.Response {
	rsp, err := c.Do(req.Request)
	if err != nil {
		log.Errorf("Failed Crawl %s %s\n!", req.RawURL, err)
		return nil
	}
	irsp := &types.Response{
		Response: rsp,
	}
	return irsp
}
