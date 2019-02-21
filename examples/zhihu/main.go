package main

import (
	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/examples/zhihu/parser"
	"github.com/tddhit/zerg/examples/zhihu/queuer"
	"github.com/tddhit/zerg/examples/zhihu/writer"
)

func main() {
	followersParser := parser.NewFollowersParser("followers")
	questionsParser := parser.NewQuestionsParser("questions")
	user2QuestionsWriter := writer.NewFileWriter("user2questions", "user2questions.txt")
	engine := engine.NewEngine(engine.Option{LogLevel: log.INFO})
	engine.AddParser(followersParser).AddParser(questionsParser)
	engine.AddWriter(user2QuestionsWriter)
	engine.SetSchedulerPolicy(queuer.NewDelayQueuer())
	engine.AddSeed("https://www.zhihu.com/people/excited-vczh/followers?page=1", "followers")
	engine.AddSeed("https://www.zhihu.com/people/excited-vczh/following?page=1", "followers")
	engine.AddSeed("https://www.zhihu.com/people/excited-vczh/following/questions?page=1", "questions")
	engine.Go()
}
