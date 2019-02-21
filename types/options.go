package types

import (
	"io"
	"net/http"
)

var defaultRequestOptions = requestOptions{}

type requestOptions struct {
	method  string
	proxy   string
	crawler string
	body    io.Reader
	header  http.Header
}

type RequestOption func(*requestOptions)

func WithMethod(method string) RequestOption {
	return func(o *requestOptions) {
		o.method = method
	}
}

func WithProxy(proxy string) RequestOption {
	return func(o *requestOptions) {
		o.proxy = proxy
	}
}

func WithCrawler(crawler string) RequestOption {
	return func(o *requestOptions) {
		o.crawler = crawler
	}
}

func WithBody(body io.Reader) RequestOption {
	return func(o *requestOptions) {
		o.body = body
	}
}

func WithHeader(header http.Header) RequestOption {
	return func(o *requestOptions) {
		o.header = header
	}
}
