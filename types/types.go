package types

import (
	"net/http"
)

type Request struct {
	*http.Request
	RawURL  string
	Parser  string
	Proxy   string
	Crawler string
}

type Response struct {
	*http.Response
	RawURL string
	Parser string
}

type Item struct {
	Dict   map[string]interface{}
	RawURL string
	Writer string
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
		Request: req,
		RawURL:  url,
		Parser:  parser,
		Proxy:   opt.proxy,
		Crawler: opt.crawler,
	}
	return ireq, nil
}

func NewItem(writer string) *Item {
	i := &Item{
		Dict:   make(map[string]interface{}),
		Writer: writer,
	}
	return i
}
