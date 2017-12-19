package engine

import (
	"bufio"
	"os"

	"github.com/tddhit/zerg/downloader"
	"github.com/tddhit/zerg/pipeline"
	"github.com/tddhit/zerg/scheduler"
	"github.com/tddhit/zerg/spider"
	"github.com/tddhit/zerg/types"
	"github.com/tddhit/zerg/util"
)

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

func NewEngine(option util.Option) *Engine {
	e := &Engine{
		reqToSchedulerChan:    make(chan *types.Request, 100),
		reqFromSchedulerChan:  make(chan *types.Request, 100),
		reqFromSpiderChan:     make(chan *types.Request, 100),
		rspToSpiderChan:       make(chan *types.Response, 100),
		itemFromSpiderChan:    make(chan *types.Item, 100),
		reqToDownloaderChan:   make(chan *types.Request, 100),
		rspFromDownloaderChan: make(chan *types.Response, 100),
		itemToPipelineChan:    make(chan *types.Item, 100),
	}
	e.scheduler = scheduler.NewScheduler(e.reqToSchedulerChan, e.reqFromSchedulerChan)
	e.downloader = downloader.NewDownloader(e.reqToDownloaderChan, e.rspFromDownloaderChan)
	e.pipeline = pipeline.NewPipeline(e.itemToPipelineChan)
	e.spider = spider.NewSpider(e.reqFromSpiderChan, e.itemFromSpiderChan, e.rspToSpiderChan)
	util.InitLogger(option)
	return e
}

func (e *Engine) AddParser(parser spider.Parser) *Engine {
	if parser != nil {
		e.spider.AddParser(parser)
	} else {
		util.LogFatalf("parser is nil!")
	}
	return e
}

func (e *Engine) SetSchedulerPolicy(q scheduler.Queuer) {
	if q != nil {
		e.scheduler.SetQueuer(q)
	} else {
		util.LogFatalf("queuer is nil!")
	}
}

func (e *Engine) AddWriter(writer pipeline.Writer) *Engine {
	if writer != nil {
		e.pipeline.AddWriter(writer)
	} else {
		util.LogFatalf("writer is nil!")
	}
	return e
}

func (e *Engine) AddSeed(url, parser string) *Engine {
	e.spider.AddSeed(url, parser)
	return e
}

func (e *Engine) AddSeedByFile(path, parser string) *Engine {
	go e.addSeedByFile(path, parser)
	return e
}

func (e *Engine) addSeedByFile(path, parser string) *Engine {
	file, err := os.Open(path)
	if err != nil {
		util.LogFatalf("AddSeedByFile fail:", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Text()
		util.LogDebugf(data)
		e.spider.AddSeed(data, parser)
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
			util.LogDebugf("spider -> engine, req:%s", req.RawURL)
			select {
			case e.reqToSchedulerChan <- req:
				util.LogDebugf("engine -> scheduler, req:%s", req.RawURL)
			default:
				util.LogWarnf("engine -> scheduler, chan is full, discard %s!", req.RawURL)
			}
		case item := <-e.itemFromSpiderChan:
			util.LogDebugf("spider -> engine, item:%s", item.RawURL)
			select {
			case e.itemToPipelineChan <- item:
				util.LogDebugf("engine -> pipeline, item:%s", item.RawURL)
			default:
				util.LogWarnf("engine -> pipeline, chan is full, discard %s!", item.RawURL)
			}
		case req := <-e.reqFromSchedulerChan:
			util.LogDebugf("scheduler -> engine, req:%s", req.RawURL)
			select {
			case e.reqToDownloaderChan <- req:
				util.LogDebugf("engine -> downloader, req:%s", req.RawURL)
			default:
				select {
				case e.reqToSchedulerChan <- req:
					util.LogWarnf("engine -> downloader, chan is full, engine -> scheduler %s !", req.RawURL)
				default:
					util.LogWarnf("engine -> downloader && engine -> scheduler, chan is full, discard %s !", req.RawURL)
				}
			}
		case rsp := <-e.rspFromDownloaderChan:
			if rsp != nil {
				util.LogDebugf("downloader -> engine, rsp:%s(%s)", rsp.RawURL, rsp.Status)
				select {
				case e.rspToSpiderChan <- rsp:
					util.LogDebugf("engine -> spider, rsp:%s", rsp.RawURL)
				default:
					util.LogWarnf("engine -> spider, chan is full, discard %s!", rsp.RawURL)
				}
			}
		}
	}
}
