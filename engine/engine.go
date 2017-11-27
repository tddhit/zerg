package engine

import (
	"log"

	"github.com/tddhit/zerg/downloader"
	"github.com/tddhit/zerg/pipeline"
	"github.com/tddhit/zerg/scheduler"
	"github.com/tddhit/zerg/spider"
	"github.com/tddhit/zerg/types"
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

func NewEngine() *Engine {
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
		log.Printf("Warning: spider %s is already exist!", spider.Name)
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
			select {
			case req := <-e.reqFromSpiderChans[name]:
				select {
				case e.reqToSchedulerChan <- req:
				default:
					log.Println("Warning: req -> scheduler is full, discard!")
				}
			case item := <-e.itemFromSpiderChans[name]:
				select {
				case e.itemToPipelineChan <- item:
				default:
					log.Println("Warning: item -> pipeline is full, discard!")
				}
			}
		}()
	}
	for {
		select {
		case req := <-e.reqFromSchedulerChan:
			select {
			case e.reqToDownloaderChan <- req:
			default:
				log.Println("Warning: req -> downloader is full, discard!")
			}
		case rsp := <-e.rspFromDownloaderChan:
			select {
			case e.rspToSpiderChans[rsp.Spider] <- rsp:
			default:
				log.Println("Warning: rsp -> spider is full, discard!")
			}
		}
	}
}
