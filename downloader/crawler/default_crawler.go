package crawler

import (
	"net/http"

	"github.com/tddhit/tools/log"

	"github.com/tddhit/zerg/types"
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
		log.Errorf("Failed Crawl %s %s\n!", req.RawURL, err)
		return nil
	}
	irsp := &types.Response{
		Response: rsp,
	}
	return irsp
}
