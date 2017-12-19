package types

import (
	"net/http"
)

type Request struct {
	*http.Request
	RawURL string
	Parser string
}

type Response struct {
	*http.Response
	RawURL string
	Parser string
}

type Item struct {
	Dict   map[string]string
	RawURL string
	Writer string
}

func NewRequest(url, parser string) (*Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "close")
	ireq := &Request{
		Request: req,
		RawURL:  url,
		Parser:  parser,
	}
	return ireq, nil
}

func NewItem(writer string) *Item {
	i := &Item{
		Dict:   make(map[string]string),
		Writer: writer,
	}
	return i
}
