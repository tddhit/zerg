package engine

import (
	"github.com/tddhit/zerg/downloader"
	"github.com/tddhit/zerg/pipeline"
	"github.com/tddhit/zerg/scheduler"
	"github.com/tddhit/zerg/spider"
	"github.com/tddhit/zerg/types"
	"github.com/tddhit/zerg/util"
)

type Engine struct {
	spiders    map[string]*spider.Spider
	scheduler  *scheduler.Scheduler
	downloader *downloader.Downloader
	pipeline   *pipeline.Pipeline

	//与scheduler通信
	reqToSchedulerChan   chan *types.Request
	reqFromSchedulerChan chan *types.Request

	//与spider通信
	reqFromSpiderChans  map[string]chan *types.Request
	rspToSpiderChans    map[string]chan *types.Response
	itemFromSpiderChans map[string]chan *types.Item

	//与downloader通信
	reqToDownloaderChan   chan *types.Request
	rspFromDownloaderChan chan *types.Response

	//与pipeline通信
	itemToPipelineChan chan *types.Item
}

func NewEngine(option util.Option) *Engine {
	e := &Engine{
		spiders:               make(map[string]*spider.Spider),
		reqToSchedulerChan:    make(chan *types.Request, 100),
		reqFromSchedulerChan:  make(chan *types.Request, 100),
		reqFromSpiderChans:    make(map[string]chan *types.Request),
		rspToSpiderChans:      make(map[string]chan *types.Response),
		itemFromSpiderChans:   make(map[string]chan *types.Item),
		reqToDownloaderChan:   make(chan *types.Request, 100),
		rspFromDownloaderChan: make(chan *types.Response, 100),
		itemToPipelineChan:    make(chan *types.Item, 100),
	}
	e.scheduler = scheduler.NewScheduler(e.reqToSchedulerChan, e.reqFromSchedulerChan)
	e.downloader = downloader.NewDownloader(e.reqToDownloaderChan, e.rspFromDownloaderChan)
	e.pipeline = pipeline.NewPipeline(e.itemToPipelineChan)
	util.InitLogger(option)
	return e
}

func (e *Engine) AddSpider(spider *spider.Spider) *Engine {
	if _, ok := e.spiders[spider.Name]; !ok {
		e.reqFromSpiderChans[spider.Name] = make(chan *types.Request, 100)
		e.itemFromSpiderChans[spider.Name] = make(chan *types.Item, 100)
		e.rspToSpiderChans[spider.Name] = make(chan *types.Response, 100)
		spider.SetupChan(e.reqFromSpiderChans[spider.Name],
			e.itemFromSpiderChans[spider.Name], e.rspToSpiderChans[spider.Name])
		e.spiders[spider.Name] = spider
	} else {
		util.LogWarn("spider[%s] is already exist!", spider.Name)
	}
	return e
}

/*
func (e *Engine) SetSchedulerPolicy(q Queuer) {
	e.scheduler = scheduler.SetQueuer(q)
}

func (e *Engine) AddDownloaderPolicy(c Crawler) {
	e.downloader.AddCrawler(c)
}

func (e *Engine) SetPipelinePolicy(w Writer) {
	e.pipeline.Add(w)
}
*/

func (e *Engine) Start() {
	e.scheduler.Go()
	e.downloader.Go()
	e.pipeline.Go()
	for name, spider := range e.spiders {
		spider.Go()
		go func() {
			for {
				select {
				case req := <-e.reqFromSpiderChans[name]:
					util.LogDebug("spider[%s] -> engine, req:%s", name, req.RawURL)
					select {
					case e.reqToSchedulerChan <- req:
						util.LogDebug("engine -> scheduler, req:%s", req.RawURL)
					default:
						util.LogWarn("engine -> scheduler, chan is full, discard %s!", req.RawURL)
					}
				case item := <-e.itemFromSpiderChans[name]:
					util.LogDebug("spider[%s] -> engine, item:%s", name, item.RawURL)
					select {
					case e.itemToPipelineChan <- item:
						util.LogDebug("engine -> pipeline, item:%s", item.RawURL)
					default:
						util.LogWarn("engine -> pipeline, chan is full, discard %s!", item.RawURL)
					}
				}
			}
		}()
	}
	for {
		select {
		case req := <-e.reqFromSchedulerChan:
			util.LogDebug("scheduler -> engine, req:%s", req.RawURL)
			select {
			case e.reqToDownloaderChan <- req:
				util.LogDebug("engine -> downloader, req:%s", req.RawURL)
			default:
				util.LogWarn("engine -> downloader, chan is full, discard %s!")
			}
		case rsp := <-e.rspFromDownloaderChan:
			if rsp != nil {
				util.LogDebug("downloader -> engine, rsp:%s(%s)", rsp.RawURL, rsp.Status)
				select {
				case e.rspToSpiderChans[rsp.Spider] <- rsp:
					util.LogDebug("engine -> spider[%s], rsp:%s", rsp.Spider, rsp.RawURL)
				default:
					util.LogWarn("engine -> spider[%s], chan is full, discard %s!", rsp.Spider, rsp.RawURL)
				}
			}
		}
	}
}
