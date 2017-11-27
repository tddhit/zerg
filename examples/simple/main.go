package main

import (
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/spider"
	"github.com/tddhit/zerg/types"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(rsp *types.Response) (*types.Item, []*types.Request) {
	item := types.NewItem()
	reqs := make([]*types.Request, 0)
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	title := doc.Find(".entry-header h1").Text()
	if title != "" {
		item.Dict["url"] = rsp.RawURL
		item.Dict["title"] = title
	} else {
		doc.Find("#archive .post-thumb a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := types.NewRequest(href)
			reqs = append(reqs, req)
		})
		doc.Find(".next.page-numbers").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := types.NewRequest(href)
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
