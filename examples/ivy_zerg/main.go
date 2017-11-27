package main

import (
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/spider"
	"github.com/tddhit/zerg/util"

	"github.com/tddhit/zerg/examples/ivy_zerg/spider/cnblogs"
	"github.com/tddhit/zerg/examples/ivy_zerg/spider/jobbole"
)

func main() {

	jobboleSpider := spider.NewSpider("jobbole", jobbole.NewParser())
	jobboleSpider.AddSeed("http://blog.jobbole.com/all-posts/")
	//jobboleSppider.AssociateWriter(pipeline.NewConsoleWriter())

	cnblogsSpider := spider.NewSpider("cnblogs", cnblogs.NewParser())
	cnblogsSpider.AddSeed("http://www.cnblogs.com")
	//cnblogsSpider.AssociateWriter(pipeline.NewFileWriter())

	engine := engine.NewEngine(util.Option{LogLevel: util.INFO})
	engine.AddSpider(jobboleSpider).AddSpider(cnblogsSpider)
	//engine.SetSchedulerPolicy(scheduler.NewQueuer())
	//engine.AddDownloaderPolicy(downloader.NewCrawler())
	engine.Start()
}
