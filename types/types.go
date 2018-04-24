package types

import (
	"net/http"
)

type Request struct {
	*http.Request
	RawURL string
	Parser string
	Proxy  string
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

func NewRequest(url, parser, proxy string, header http.Header) (*Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v1 := range header {
		for _, v2 := range v1 {
			req.Header.Add(k, v2)
		}
	}
	req.Header.Set("Connection", "close")
	ireq := &Request{
		Request: req,
		RawURL:  url,
		Parser:  parser,
		Proxy:   proxy,
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
