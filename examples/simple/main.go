package main

import (
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/util"

	"github.com/tddhit/zerg/examples/simple/parser"
	"github.com/tddhit/zerg/examples/simple/writer"
)

func main() {
	jobboleParser := parser.NewJobboleParser("jobbole")
	cnblogsParser := parser.NewCnblogsParser("cnblogs")

	cnblogsWriter := writer.NewFileWriter("cnblogs", "cnblogs.txt")

	engine := engine.NewEngine(util.Option{LogLevel: util.INFO})
	engine.AddParser(cnblogsParser).AddParser(jobboleParser)
	engine.AddWriter(cnblogsWriter)
	engine.AddSeed("http://blog.jobbole.com/all-posts/", "jobbole")
	engine.AddSeed("http://www.cnblogs.com", "cnblogs")
	engine.Go()
}
