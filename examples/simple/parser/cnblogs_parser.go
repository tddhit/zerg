package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/zerg/types"
)

const ROOTURL = "https://www.cnblogs.com/"

type CnblogsParser struct {
	name string
}

func NewCnblogsParser(name string) *CnblogsParser {
	return &CnblogsParser{
		name: name,
	}
}

func (p *CnblogsParser) Name() string {
	return p.name
}

func (p *CnblogsParser) Parse(rsp *types.Response) (*types.Item, []*types.Request) {
	item := types.NewItem("cnblogs")
	reqs := make([]*types.Request, 0)
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	title := doc.Find("#cb_post_title_url").Text()
	if title != "" {
		item.Dict["url"] = rsp.RawURL
		item.Dict["title"] = title
	} else {
		doc.Find("#post_list .post_item_body h3 a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := types.NewRequest(href, rsp.Parser)
			reqs = append(reqs, req)
		})
		doc.Find("#pager_bottom .pager :last-child").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := types.NewRequest(ROOTURL+href, rsp.Parser)
			reqs = append(reqs, req)
		})
	}
	return item, reqs
}
