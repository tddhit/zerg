package types

import (
	"net/http"
)

type Request struct {
	*http.Request
	RawURL string
	Spider string
}

type Response struct {
	*http.Response
	RawURL string
	Spider string
}

type Item struct {
	//associatedWriter Writer
	Dict   map[string]string
	RawURL string
	Spider string
}

func NewRequest(url, spiderName string) (*Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	ireq := &Request{
		Request: req,
		RawURL:  url,
		Spider:  spiderName,
	}
	return ireq, nil
}

func NewItem() *Item {
	i := &Item{
		Dict: make(map[string]string),
	}
	return i
}
