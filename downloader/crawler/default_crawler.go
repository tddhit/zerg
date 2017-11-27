package crawler

import (
	"log"
	"net/http"

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
		log.Println("Failed Crawl %s %d\n!", req.RawURL, rsp.Status)
		return nil
	}
	irsp := &types.Response{
		Response: rsp,
	}
	return irsp
}
