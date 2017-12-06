package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/zerg/types"
)

const ROOTURL = "https://movie.douban.com"

type DoubanParser struct {
	name string
}

func NewDoubanParser(name string) *DoubanParser {
	return &DoubanParser{
		name: name,
	}
}

func (p *DoubanParser) Name() string {
	return p.name
}

func (p *DoubanParser) Parse(rsp *types.Response) (*types.Item, []*types.Request) {
	item := types.NewItem("douban")
	reqs := make([]*types.Request, 0)
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	title := doc.Find("h1 span").Text()
	if title != "" {
		content := doc.Find(".main-bd p").Text()
		content = strings.Join(strings.Fields(content), " ")
		item.Dict["url"] = rsp.RawURL
		item.Dict["title"] = title
		item.Dict["content"] = content
	} else {
		doc.Find(".main-bd h2 a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := types.NewRequest(href, rsp.Parser)
			reqs = append(reqs, req)
		})
		doc.Find(".next a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := types.NewRequest(ROOTURL+href, rsp.Parser)
			reqs = append(reqs, req)
		})
	}
	return item, reqs
}
