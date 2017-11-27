package engine

import (
	"log"
	"net/http"
	"sync"

	"github.com/tddhit/zerg/downloader"
	"github.com/tddhit/zerg/pipeline"
	"github.com/tddhit/zerg/scheduler"
	"github.com/tddhit/zerg/spider"
	"github.com/tddhit/zerg/util"
)

type Request struct {
	*http.Request
}

type Response struct {
	*http.Response
}

type Item struct {
	associatedWriter Writer
	Dict             map[string]string
}

type Engine struct {
	mutex      sync.Mutex
	spiders    map[string]*spider.Spider
	scheduler  scheduler.Scheduler
	downloader downloader.Downloader
	pipeline   pipeline.Pipeline

	//与scheduler通信
	reqToSchedulerChan   chan *Request
	reqFromSchedulerChan chan *Request

	//与spider通信
	reqFromSpiderChan  chan *Request
	rspToSpiderChan    chan *Response
	itemFromSpiderChan chan *Item

	//与downloader通信
	reqToDownloaderChan   chan *Request
	rspFromDownloaderChan chan *Response

	//与pipeline通信
	itemToPipelineChan chan *Item
}

func NewEngine() *Engine {
	e := &Engine{
		spiders:               make(map[string]*spider.Spider),
		reqToSchedulerChan:    make(chan *Request, 10),
		reqFromSchedulerChan:  make(chan *Request, 10),
		reqFromSpiderChan:     make(chan *Request, 10),
		rspToSpiderChan:       make(chan *Response, 10),
		itemFromSpiderChan:    make(chan *Item, 10),
		reqToDownloaderChan:   make(chan *Request, 10),
		rspFromDownloaderChan: make(chan *Response, 10),
		itemToPipelineChan:    make(chan *Item, 10),
	}
	e.scheduler = scheduler.NewScheduler(reqToSchedulerChan, reqFromSchedulerChan)
	e.downloader = downloader.NewDownloader(reqToDownloaderChan, rspToEngineChan)
	e.pipeline = pipeline.NewPipeline(itemFromEngineChan)
	return e
}

func (e *Engine) AddSpider(spider spider.Spider) *Engine {
	if sp, ok := e.spiders[spider.Name]; !ok {
		e.spiders[spider.Name] = spider
		spider.SetupChan(reqFromSpiderChan, itemFromSpiderChan, rspToSpiderChan)
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
	scheduler.Go()
	downloader.Go()
	pipeline.Go()
	for _, spider := range spiders {
		spider.Go()
	}
	for {
		select {
		case req := <-reqFromSpiderChan:
			select {
			case reqToSchedulerChan <- req:
			default:
				log.Println("Warning: req -> scheduler is full, discard!")
			}
		case req := <-reqFromSchedulerChan:
			select {
			case reqToDownloaderChan <- req:
				log.Println("Warning: req -> downloader is full, discard!")
			default:
			}
		case rsp := <-rspFromDownloaderChan:
			select {
			case rspToSpiderChan <- rsp:
			default:
				log.Println("Warning: rsp -> spider is full, discard!")
			}
		case item := <-itemFromSpiderChan:
			select {
			case itemToPipelineChan <- item:
			default:
				log.Println("Warning: item -> pipeline is full, discard!")
			}
		}
	}
}
