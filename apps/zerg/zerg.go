package main

import (
	"flag"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/apps/zerg/parser"
	"github.com/tddhit/zerg/apps/zerg/writer"
	"github.com/tddhit/zerg/engine"
)

var confPath string

func init() {
	flag.StringVar(&confPath, "conf", "zerg.yml", "config file")
	flag.Parse()
}

type program struct {
	engine *engine.Engine
	parser map[string]*parser.Parser
	writer map[string]*writer.FileWriter
}

func newProgram(conf *Conf) *program {
	p := &program{
		engine: engine.NewEngine(engine.Option{LogLevel: log.DEBUG}),
		parser: make(map[string]*parser.Parser),
		writer: make(map[string]*writer.FileWriter),
	}
	for k, v := range conf.Parser {
		p.parser[k] = parser.NewParser(k, v.Type, v.CssSelector, v.Writer, v.Parser)
		p.engine.AddParser(p.parser[k])
	}
	for k, v := range conf.Writer {
		p.writer[k] = writer.NewFileWriter(k, v)
		p.engine.AddWriter(p.writer[k])
	}
	p.engine.AddSeedByFile(conf.Seed.File, conf.Seed.Parser)
	return p
}

func (p *program) Go() {
	p.engine.Go()
	select {}
}

func main() {
	conf, err := NewConf(confPath)
	if err != nil {
		log.Fatal("parser config fail.")
	}
	p := newProgram(conf)
	p.Go()
}
