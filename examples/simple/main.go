package main

import (
	"log"

	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/spider"
)

type Parser struct {
	reg *regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(rsp *engine.Response) (*engine.Item, []*engine.Request) {
	item := &engine.Item{}
	reqs := make([]*engine.Request, 0)
	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	title := doc.Find(".entry-header h1").Text()
	if title != "" {
		item.Dict["title"] = title
	} else {
		doc.Find("#archive .post-thumb a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := http.NewRequest("GET", href, nil)
			ireq := &engine.Request{
				req,
				RawURL: href,
			}
			reqs = append(reqs, req)
		})
	}
	return item, reqs
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	engine := engine.NewEngine()

	jobboleSpider := spider.NewSpider("jobbole", NewParser())
	jobboleSpider.AddSeed("http://blog.jobbole.com/all-posts/")

	engine.AddSpider(jobboleSpider)
	engine.Start()
}
