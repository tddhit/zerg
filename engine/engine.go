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
	reqFromSpiderChan  chan *types.Request
	rspToSpiderChan    chan *types.Response
	itemFromSpiderChan chan *types.Item

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
	return e
}

func (e *Engine) AddSpider(spider *spider.Spider) *Engine {
	if _, ok := e.spiders[spider.Name]; !ok {
		e.spiders[spider.Name] = spider
		spider.SetupChan(e.reqFromSpiderChan, e.itemFromSpiderChan, e.rspToSpiderChan)
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
	for _, spider := range e.spiders {
		spider.Go()
	}
	for {
		select {
		case req := <-e.reqFromSpiderChan:
			select {
			case e.reqToSchedulerChan <- req:
			default:
				log.Println("Warning: req -> scheduler is full, discard!")
			}
		case req := <-e.reqFromSchedulerChan:
			select {
			case e.reqToDownloaderChan <- req:
			default:
				log.Println("Warning: req -> downloader is full, discard!")
			}
		case rsp := <-e.rspFromDownloaderChan:
			select {
			case e.rspToSpiderChan <- rsp:
			default:
				log.Println("Warning: rsp -> spider is full, discard!")
			}
		case item := <-e.itemFromSpiderChan:
			select {
			case e.itemToPipelineChan <- item:
			default:
				log.Println("Warning: item -> pipeline is full, discard!")
			}
		}
	}
}
