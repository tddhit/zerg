package zerg

import (
	"fmt"
	"net/http"
)

type Request struct {
	*http.Request
	RawURL   string
	Parser   string
	Proxy    string
	Crawler  string
	ID       string
	Metadata map[string]string
}

type Response struct {
	*http.Response
	*Request
	Err error
}

type Item struct {
	Dict     map[string]interface{}
	RawURL   string
	Writers  []string
	Metadata map[string]string
}

func NewRequest(url, parser string, opts ...RequestOption) (*Request, error) {
	opt := defaultRequestOptions
	for _, o := range opts {
		o(&opt)
	}
	req, err := http.NewRequest(opt.method, url, opt.body)
	if err != nil {
		return nil, err
	}
	for k, v1 := range opt.header {
		for _, v2 := range v1 {
			req.Header.Add(k, v2)
		}
	}
	req.Header.Set("Connection", "close")
	ireq := &Request{
		Request:  req,
		RawURL:   url,
		Parser:   parser,
		Proxy:    opt.proxy,
		Crawler:  opt.crawler,
		Metadata: opt.metadata,
		ID:       opt.id,
	}
	if ireq.Crawler == "" {
		ireq.Crawler = "DEFAULT_CRAWLER"
	}
	return ireq, nil
}

func (r *Request) String() string {
	return fmt.Sprintf("id[%s] %s", r.ID, r.RawURL)
}

func NewItem(writers ...string) *Item {
	i := &Item{
		Dict:     make(map[string]interface{}),
		Writers:  writers,
		Metadata: make(map[string]string),
	}
	if len(i.Writers) == 0 {
		i.Writers = append(i.Writers, "DEFAULT_WRITER")
	}
	return i
}
