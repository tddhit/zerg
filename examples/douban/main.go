package main

import (
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/tools/log"

	"github.com/tddhit/zerg/examples/douban/parser"
	"github.com/tddhit/zerg/examples/douban/queuer"
	"github.com/tddhit/zerg/examples/douban/writer"
)

func main() {
	doubanParser := parser.NewDoubanParser("douban")
	doubanWriter := writer.NewFileWriter("douban", "douban.txt")
	sleepQueuer := queuer.NewDefaultQueuer()

	engine := engine.NewEngine(util.Option{LogLevel: util.DEBUG})
	engine.AddParser(doubanParser)
	engine.AddWriter(doubanWriter)
	engine.SetSchedulerPolicy(sleepQueuer)
	engine.AddSeed("https://movie.douban.com/review/best/", "douban")
	engine.Go()
}
