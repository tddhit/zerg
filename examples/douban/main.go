package main

import (
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/util"

	"github.com/tddhit/zerg/examples/douban/parser"
	"github.com/tddhit/zerg/examples/douban/writer"
)

func main() {
	doubanParser := parser.NewDoubanParser("douban")
	doubanWriter := writer.NewFileWriter("douban", "douban.txt")

	engine := engine.NewEngine(util.Option{LogLevel: util.INFO})
	engine.AddParser(doubanParser)
	engine.AddWriter(doubanWriter)
	engine.AddSeed("https://movie.douban.com/review/best/", "douban")
	engine.Go()
}
