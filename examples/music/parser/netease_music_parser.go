package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/zerg/types"
)

const ROOTURL = "http://music.163.com/album?"

type NeteaseMusicParser struct {
	name string
}

func NewNeteaseMusicParser(name string) *NeteaseMusicParser {
	return &NeteaseMusicParser{
		name: name,
	}
}

func (p *NeteaseMusicParser) Name() string {
	return p.name
}

func (p *NeteaseMusicParser) Parse(rsp *types.Response) ([]*types.Item, []*types.Request) {
	var items []*types.Item
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	if doc == nil {
		return nil, nil
	}
	artist := doc.Find(".intr span").Text()
	doc.Find(".f-hide li a").Each(func(i int, contentSelection *goquery.Selection) {
		musicName := contentSelection.Text()
		item := types.NewItem("NeteaseMusicWriter")
		item.Dict["Artist"] = artist
		item.Dict["MusicName"] = musicName
		items = append(items, item)
	})
	return items, nil
}
