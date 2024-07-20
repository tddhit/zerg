package zerg

import (
	"io"
	"net/http"
)

type zergOptions struct {
	logPath  string
	logLevel int
	crawlers []Crawler
	writers  []Writer
	parsers  []Parser
	queuer   Queuer
}

type ZergOption func(*zergOptions)

var defaultZergOptions = zergOptions{}

func WithLogPath(path string) ZergOption {
	return func(o *zergOptions) {
		o.logPath = path
	}
}

func WithLogLevel(level int) ZergOption {
	return func(o *zergOptions) {
		o.logLevel = level
	}
}

func WithCrawler(c Crawler) ZergOption {
	return func(o *zergOptions) {
		o.crawlers = append(o.crawlers, c)
	}
}

func WithParser(p Parser) ZergOption {
	return func(o *zergOptions) {
		o.parsers = append(o.parsers, p)
	}
}

func WithWriter(w Writer) ZergOption {
	return func(o *zergOptions) {
		o.writers = append(o.writers, w)
	}
}

func WithQueuer(q Queuer) ZergOption {
	return func(o *zergOptions) {
		o.queuer = q
	}
}

type requestOptions struct {
	id       string
	method   string
	proxy    string
	crawler  string
	body     io.Reader
	header   http.Header
	metadata map[string]string
}

type RequestOption func(*requestOptions)

var defaultRequestOptions = requestOptions{}

func WithRequestID(id string) RequestOption {
	return func(o *requestOptions) {
		o.id = id
	}
}

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

func WithRequestCrawler(crawler string) RequestOption {
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

func WithMetadata(meta map[string]string) RequestOption {
	return func(o *requestOptions) {
		o.metadata = meta
	}
}
