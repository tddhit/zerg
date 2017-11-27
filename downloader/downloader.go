package downloader

type Crawler interface {
	Crawl(req *Request) *Response
}

type Downloader struct {
	Crawler
	reqFromEngineChan <-chan *Request
	rspToEngineChan   chan<- *Response
}

func NewDownloader(reqFromEngineChan <-chan *Request, rspToEngineChan chan<- *Response) {
	d := &Downloader{
		Crawler:           crawler.NewHTTPCrawler(),
		reqFromEngineChan: reqFromEngineChan,
		rspToEngineChan:   rspToEngineChan,
	}
}

func (d *Downloader) Go() {
	go func() {
		for {
			req := <-reqFromEngineChan
			rsp := d.Crawl(req)
			rspToEngineChan <- rsp
		}
	}()
}
