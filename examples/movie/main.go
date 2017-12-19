package main

import (
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/examples/movie/parser"
	"github.com/tddhit/zerg/examples/movie/writer"
	"github.com/tddhit/zerg/util"
)

func main() {
	baidunewsParser := parser.NewBaiduNewsParser("baidunews")
	fullParser := parser.NewFullParser("full")
	fullWriter := writer.NewFileWriter("full", "data/full.txt")

	engine := engine.NewEngine(util.Option{LogLevel: util.WARNING})
	engine.AddParser(baidunewsParser)
	engine.AddParser(fullParser).AddWriter(fullWriter)
	engine.AddSeedByFile("data/baidunews_movie.txt", "baidunews")
	engine.Go()
}
