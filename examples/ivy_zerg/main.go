package main

import (
	"log"

	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/spider"

	//"./downloader"
	//"./pipeline"
	//"./scheduler"
	"github.com/tddhit/zerg/examples/ivy_zerg/spider/cnblogs"
	"github.com/tddhit/zerg/examples/ivy_zerg/spider/jobbole"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	engine := engine.NewEngine()

	jobboleSpider := spider.NewSpider("jobbole", jobbole.NewParser())
	jobboleSpider.AddSeed("http://blog.jobbole.com/all-posts/")
	//jobboleSppider.AssociateWriter(pipeline.NewConsoleWriter())

	cnblogsSpider := spider.NewSpider("cnblogs", cnblogs.NewParser())
	cnblogsSpider.AddSeed("http://www.cnblogs.com")
	//cnblogsSpider.AssociateWriter(pipeline.NewFileWriter())

	engine.AddSpider(jobboleSpider).AddSpider(cnblogsSpider)
	//engine.SetSchedulerPolicy(scheduler.NewQueuer())
	//engine.AddDownloaderPolicy(downloader.NewCrawler())
	engine.Start()
}
