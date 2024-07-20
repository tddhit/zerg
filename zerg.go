package zerg

import (
	"bufio"
	"errors"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync/atomic"
	"time"

	"github.com/tddhit/tools/log"
)

type Zerg struct {
	opt        *zergOptions
	spider     *spider
	scheduler  *scheduler
	downloader *downloader
	pipeline   *pipeline

	reqToSchedulerC   chan *Request
	reqFromSchedulerC chan *Request

	reqFromSpiderC  chan *Request
	rspToSpiderC    chan *Response
	itemFromSpiderC chan *Item

	reqToDownloaderC   chan *Request
	rspFromDownloaderC chan *Response

	itemToPipelineC chan *Item

	closeC   chan struct{}
	stopFlag int32
}

func New(opts ...ZergOption) (*Zerg, error) {
	opt := defaultZergOptions
	for _, o := range opts {
		o(&opt)
	}
	if len(opt.crawlers) == 0 ||
		len(opt.parsers) == 0 ||
		len(opt.writers) == 0 ||
		opt.queuer == nil {

		return nil, errors.New("There is at least one parser/crawler/writer/queuer")
	}
	z := &Zerg{
		opt:                &opt,
		reqToSchedulerC:    make(chan *Request, 1000),
		reqFromSchedulerC:  make(chan *Request, 1000),
		reqFromSpiderC:     make(chan *Request, 1000),
		rspToSpiderC:       make(chan *Response, 1000),
		itemFromSpiderC:    make(chan *Item, 1000),
		reqToDownloaderC:   make(chan *Request, 1000),
		rspFromDownloaderC: make(chan *Response, 1000),
		itemToPipelineC:    make(chan *Item, 1000),
		closeC:             make(chan struct{}),
	}
	log.Init(opt.logPath, opt.logLevel)
	z.scheduler = newScheduler(z.reqToSchedulerC, z.reqFromSchedulerC)
	z.downloader = newDownloader(z.reqToDownloaderC, z.rspFromDownloaderC)
	z.pipeline = newPipeline(z.itemToPipelineC)
	z.spider = newSpider(z.reqFromSpiderC, z.itemFromSpiderC, z.rspToSpiderC)
	for _, c := range opt.crawlers {
		z.downloader.addCrawler(c)
	}
	for _, p := range opt.parsers {
		z.spider.addParser(p)
	}
	for _, w := range opt.writers {
		z.pipeline.addWriter(w)
	}
	z.scheduler.setQueuer(opt.queuer)
	return z, nil
}

func (z *Zerg) AddProxy(proxy string) {
	z.downloader.addProxy(proxy)
}

func (z *Zerg) RemoveProxy(proxy string) {
	z.downloader.removeProxy(proxy)
}

func (z *Zerg) AddSeed(url, parser string) error {
	return z.spider.addSeed(url, parser)
}

func (z *Zerg) AddSeedByFile(path, parser string) {
	go func() {
		file, err := os.Open(path)
		if err != nil {
			log.Fatal("AddSeedByFile fail:", err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data := scanner.Text()
			if err := z.spider.addSeed(data, parser); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Millisecond * 20)
		}
	}()
}

func (z *Zerg) AddRequest(req *Request) {
	z.spider.addRequest(req)
}

func (z *Zerg) Start() {
	go z.scheduler.start()
	go z.downloader.start()
	go z.pipeline.start()
	go z.spider.start()
	go func() {
		http.ListenAndServe(":7777", nil)
	}()
	for {
		if atomic.LoadInt32(&z.stopFlag) == 1 {
			z.prepareToStop()
		}
		select {
		case req := <-z.reqFromSpiderC:
			log.Debugf("spider -> engine, req:%s", req.String())
			z.reqToSchedulerC <- req
			log.Debugf("engine -> scheduler, req:%s", req.String())
		case item := <-z.itemFromSpiderC:
			log.Debugf("spider -> engine, item:%s", item.RawURL)
			z.itemToPipelineC <- item
			log.Debugf("engine -> pipeline, item:%s", item.RawURL)
		case req, ok := <-z.reqFromSchedulerC:
			if !ok {
				z.reqFromSchedulerC = nil
				atomic.StoreInt32(&z.stopFlag, 1)
				continue
			}
			log.Debugf("scheduler -> engine, req:%s", req.String())
			z.reqToDownloaderC <- req
			log.Debugf("engine -> downloader, req:%s", req.String())
		case rsp, ok := <-z.rspFromDownloaderC:
			if !ok {
				z.rspFromDownloaderC = nil
				z.Stop()
				continue
			}
			if rsp == nil {
				continue
			}
			log.Debugf("downloader -> engine, req:%s", rsp.RawURL)
			if rsp.Err != nil {
				z.reqToSchedulerC <- rsp.Request
				log.Debugf("engine -> scheduler(requeue), req:%s, err:%s",
					rsp.RawURL, rsp.Err)
			} else {
				z.rspToSpiderC <- rsp
				log.Debugf("engine -> spider, rsp:%s", rsp.RawURL)
			}
		case <-z.closeC:
			goto exit
		}
	}
exit:
	z.closeChannel()
}

func (z *Zerg) prepareToStop() {
	count := len(z.reqToSchedulerC)
	count += len(z.reqFromSpiderC) + len(z.itemFromSpiderC) + len(z.rspToSpiderC)
	count += len(z.reqToDownloaderC) + len(z.rspFromDownloaderC)
	count += len(z.itemToPipelineC)
	if count == 0 {
		close(z.closeC)
	}
}

func (z *Zerg) closeChannel() {
	if z.reqFromSchedulerC != nil {
		close(z.reqFromSchedulerC)
	}
	if z.reqToSchedulerC != nil {
		close(z.reqToSchedulerC)
	}
	if z.reqFromSpiderC != nil {
		close(z.reqFromSpiderC)
	}
	if z.rspToSpiderC != nil {
		close(z.rspToSpiderC)
	}
	if z.itemFromSpiderC != nil {
		close(z.itemFromSpiderC)
	}
	if z.reqToDownloaderC != nil {
		close(z.reqToDownloaderC)
	}
	if z.rspFromDownloaderC != nil {
		close(z.rspFromDownloaderC)
	}
	if z.itemToPipelineC != nil {
		close(z.itemToPipelineC)
	}
	z.downloader.close()
}

func (z *Zerg) Stop() {
	close(z.closeC)
}

func (z *Zerg) GracefulStop() {
	close(z.reqFromSchedulerC)
}
