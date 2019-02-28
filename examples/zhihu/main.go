package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg"
	"github.com/tddhit/zerg/examples/zhihu/parser"
	"github.com/tddhit/zerg/ext/crawler"
	"github.com/tddhit/zerg/ext/queuer"
	"github.com/tddhit/zerg/ext/writer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./zhihu ip.txt")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	engine, err := zerg.New(
		zerg.WithLogLevel(log.INFO),
		zerg.WithLogPath(""),
		zerg.WithParser(parser.NewFollowersParser("followers")),
		zerg.WithParser(parser.NewQuestionsParser("questions")),
		zerg.WithWriter(writer.NewFileWriter("user2questions", "user2questions.txt")),
		zerg.WithCrawler(crawler.NewDefaultCrawler()),
		zerg.WithQueuer(queuer.NewDefaultQueuer()),
	)
	if err != nil {
		log.Fatal(err)
	}
	//engine.AddSeed("https://www.zhihu.com/people/excited-vczh/followers?page=1", "followers")
	//engine.AddSeed("https://www.zhihu.com/people/excited-vczh/following?page=1", "followers")
	//engine.AddSeed("https://www.zhihu.com/people/excited-vczh/following/questions?page=1", "questions")
	engine.AddSeed("http://www.zhihu.com/people/pjer/following?page=1", "followers")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ip := scanner.Text()
		log.Info("add")
		engine.AddProxy(ip)
	}
	engine.Start()
}
