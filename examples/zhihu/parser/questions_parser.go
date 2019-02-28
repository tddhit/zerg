package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg"
)

const ROOTURL = "https://www.zhihu.com"

type questionsParser struct {
	name string
}

func NewQuestionsParser(name string) *questionsParser {
	return &questionsParser{
		name: name,
	}
}

func (p *questionsParser) Name() string {
	return p.name
}

func (p *questionsParser) Parse(rsp *zerg.Response) ([]*zerg.Item, []*zerg.Request) {
	tokens := strings.Split(rsp.RawURL, "/")
	if len(tokens) != 7 {
		log.Error("tokens != 7")
		return nil, nil
	}
	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		log.Error(err)
		return nil, nil
	}
	var (
		items []*zerg.Item
		reqs  []*zerg.Request
	)
	doc.Find("#Profile-following .ContentItem a").Each(
		func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			item := zerg.NewItem("user2questions")
			item.Dict["user"] = tokens[4]
			item.Dict["title"] = contentSelection.Text()
			item.Dict["url"] = ROOTURL + href
			items = append(items, item)
		})
	if len(items) > 0 {
		tokens := strings.Split(rsp.RawURL, "?page=")
		if len(tokens) == 2 {
			if page, err := strconv.Atoi(tokens[1]); err == nil {
				url := fmt.Sprintf("%s?page=%d", tokens[0], page+1)
				req, _ := zerg.NewRequest(url, "questions")
				reqs = append(reqs, req)
			}
		}
	}
	return items, reqs
}
