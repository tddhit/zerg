package crawler

import (
	"net/http"

	"github.com/tddhit/zerg/types"
	"github.com/tddhit/zerg/util"
)

type HTTPCrawler struct {
	*http.Client
}

func NewHTTPCrawler() *HTTPCrawler {
	client := &http.Client{}
	crawler := &HTTPCrawler{
		Client: client,
	}
	return crawler
}

func (c *HTTPCrawler) Crawl(req *types.Request) *types.Response {
	rsp, err := c.Do(req.Request)
	if err != nil {
		util.LogError("Failed Crawl %s %s\n!", req.RawURL, err)
		return nil
	}
	irsp := &types.Response{
		Response: rsp,
		RawURL:   req.RawURL,
		Spider:   req.Spider,
	}
	return irsp
}
