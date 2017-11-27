package downloader

import (
	"net/http"
)

type HTTPCrawler struct {
	*http.Client
}

func NewHTTPCrawler() *HTTPCrawler {
	tr := &http.Transport{
	}
	client := &http.Client{
	}
	crawler := &HTTPCrawler{
		Client: 
	}
}

func (c *HTTPCrawler) Crawl() {
}
