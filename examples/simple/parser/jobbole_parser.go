package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/zerg/types"
)

type JobboleParser struct {
	name string
}

func NewJobboleParser(name string) *JobboleParser {
	return &JobboleParser{
		name: name,
	}
}

func (p *JobboleParser) Name() string {
	return p.name
}

func (p *JobboleParser) Parse(rsp *types.Response) (*types.Item, []*types.Request) {
	item := types.NewItem("")
	reqs := make([]*types.Request, 0)
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	title := doc.Find(".entry-header h1").Text()
	if title != "" {
		item.Dict["url"] = rsp.RawURL
		item.Dict["title"] = title
	} else {
		doc.Find("#archive .post-thumb a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := types.NewRequest(href, rsp.Parser)
			reqs = append(reqs, req)
		})
		doc.Find(".next.page-numbers").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := types.NewRequest(href, rsp.Parser)
			reqs = append(reqs, req)
		})
	}
	return item, reqs
}
