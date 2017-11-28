package main

import (
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/spider"
	"github.com/tddhit/zerg/util"

	"github.com/tddhit/zerg/examples/ivy_zerg/parser"
	"github.com/tddhit/zerg/examples/ivy_zerg/queuer"
	"github.com/tddhit/zerg/examples/ivy_zerg/writer"
)

func main() {

	jobboleSpider := spider.NewSpider("jobbole", parser.NewJobboleParser())
	jobboleSpider.AddSeed("http://blog.jobbole.com/all-posts/")

	cnblogsSpider := spider.NewSpider("cnblogs", parser.NewCnblogsParser())
	cnblogsSpider.AddSeed("http://www.cnblogs.com")

	engine := engine.NewEngine(util.Option{LogLevel: util.INFO})
	engine.AddSpider(jobboleSpider).AddSpider(cnblogsSpider)
	engine.AssociateWriter(jobboleSpider, writer.NewConsoleWriter())
	engine.AssociateWriter(cnblogsSpider, writer.NewFileWriter("cnblogs.txt"))
	engine.SetSchedulerPolicy(queuer.NewChanQueuer())
	//engine.AddDownloaderPolicy(downloader.NewCrawler())
	engine.Start()
}
