package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/zerg/types"
)

const ROOTURL = "http://news.baidu.com/"

type BaiduNewsParser struct {
	name string
}

func NewBaiduNewsParser(name string) *BaiduNewsParser {
	return &BaiduNewsParser{
		name: name,
	}
}

func (p *BaiduNewsParser) Name() string {
	return p.name
}

func (p *BaiduNewsParser) Parse(rsp *types.Response) ([]*types.Item, []*types.Request) {
	var reqs []*types.Request
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	doc.Find(".result .c-title a").Each(func(i int, contentSelection *goquery.Selection) {
		href, _ := contentSelection.Attr("href")
		req, _ := types.NewRequest(href, "full")
		reqs = append(reqs, req)
	})
	//doc.Find("#page .n").Each(func(i int, contentSelection *goquery.Selection) {
	//	href, _ := contentSelection.Attr("href")
	//	req, _ := types.NewRequest(ROOTURL+href, rsp.Parser)
	//	reqs = append(reqs, req)
	//})
	return nil, reqs
}
