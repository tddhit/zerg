package engine

import (
	"bufio"
	"os"
	"time"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/downloader"
	"github.com/tddhit/zerg/pipeline"
	"github.com/tddhit/zerg/scheduler"
	"github.com/tddhit/zerg/spider"
	"github.com/tddhit/zerg/types"
)

type Option struct {
	LogPath  string
	LogLevel int
}

type Engine struct {
	spider     *spider.Spider
	scheduler  *scheduler.Scheduler
	downloader *downloader.Downloader
	pipeline   *pipeline.Pipeline

	//与scheduler通信
	reqToSchedulerChan   chan *types.Request
	reqFromSchedulerChan chan *types.Request

	//与spider通信
	reqFromSpiderChan  chan *types.Request
	rspToSpiderChan    chan *types.Response
	itemFromSpiderChan chan *types.Item

	//与downloader通信
	reqToDownloaderChan   chan *types.Request
	rspFromDownloaderChan chan *types.Response

	//与pipeline通信
	itemToPipelineChan chan *types.Item
}

func NewEngine(option Option) *Engine {
	e := &Engine{
		reqToSchedulerChan:    make(chan *types.Request, 1000),
		reqFromSchedulerChan:  make(chan *types.Request, 1000),
		reqFromSpiderChan:     make(chan *types.Request, 1000),
		rspToSpiderChan:       make(chan *types.Response, 1000),
		itemFromSpiderChan:    make(chan *types.Item, 1000),
		reqToDownloaderChan:   make(chan *types.Request, 1000),
		rspFromDownloaderChan: make(chan *types.Response, 1000),
		itemToPipelineChan:    make(chan *types.Item, 1000),
	}
	e.scheduler = scheduler.NewScheduler(e.reqToSchedulerChan, e.reqFromSchedulerChan)
	e.downloader = downloader.NewDownloader(e.reqToDownloaderChan, e.rspFromDownloaderChan)
	e.pipeline = pipeline.NewPipeline(e.itemToPipelineChan)
	e.spider = spider.NewSpider(e.reqFromSpiderChan, e.itemFromSpiderChan, e.rspToSpiderChan)
	log.Init(option.LogPath, option.LogLevel)
	return e
}

func (e *Engine) AddParser(parser spider.Parser) *Engine {
	if parser != nil {
		e.spider.AddParser(parser)
	} else {
		log.Fatalf("parser is nil!")
	}
	return e
}

func (e *Engine) SetSchedulerPolicy(q scheduler.Queuer) {
	if q != nil {
		e.scheduler.SetQueuer(q)
	} else {
		log.Fatalf("queuer is nil!")
	}
}

func (e *Engine) AddWriter(writer pipeline.Writer) *Engine {
	if writer != nil {
		e.pipeline.AddWriter(writer)
	} else {
		log.Fatalf("writer is nil!")
	}
	return e
}

func (e *Engine) AddSeed(url, parser, proxy string) *Engine {
	e.spider.AddSeed(url, parser, proxy)
	return e
}

func (e *Engine) AddSeedByFile(path, parser string) *Engine {
	go e.addSeedByFile(path, parser)
	return e
}

func (e *Engine) addSeedByFile(path, parser string) *Engine {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("AddSeedByFile fail:", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Text()
		e.spider.AddSeed(data, parser, "")
		time.Sleep(time.Millisecond * 20)
	}
	return e
}

func (e *Engine) Go() {
	e.scheduler.Go()
	e.downloader.Go()
	e.pipeline.Go()
	e.spider.Go()
	for {
		select {
		case req := <-e.reqFromSpiderChan:
			log.Debugf("spider -> engine, req:%s", req.RawURL)
			select {
			case e.reqToSchedulerChan <- req:
				log.Debugf("engine -> scheduler, req:%s", req.RawURL)
			default:
				log.Warnf("engine -> scheduler, chan is full, discard %s!", req.RawURL)
			}
		case item := <-e.itemFromSpiderChan:
			log.Debugf("spider -> engine, item:%s", item.RawURL)
			select {
			case e.itemToPipelineChan <- item:
				log.Debugf("engine -> pipeline, item:%s", item.RawURL)
			default:
				log.Warnf("engine -> pipeline, chan is full, discard %s!", item.RawURL)
			}
		case req := <-e.reqFromSchedulerChan:
			log.Debugf("scheduler -> engine, req:%s", req.RawURL)
			select {
			case e.reqToDownloaderChan <- req:
				log.Debugf("engine -> downloader, req:%s", req.RawURL)
			default:
				select {
				case e.reqToSchedulerChan <- req:
					log.Warnf("engine -> downloader, chan is full, engine -> scheduler %s !", req.RawURL)
				default:
					log.Warnf("engine -> downloader && engine -> scheduler, chan is full, discard %s !", req.RawURL)
				}
			}
		case rsp := <-e.rspFromDownloaderChan:
			if rsp != nil {
				log.Debugf("downloader -> engine, rsp:%s(%s)", rsp.RawURL, rsp.Status)
				select {
				case e.rspToSpiderChan <- rsp:
					log.Debugf("engine -> spider, rsp:%s", rsp.RawURL)
				default:
					log.Warnf("engine -> spider, chan is full, discard %s!", rsp.RawURL)
				}
			}
		}
	}
}
