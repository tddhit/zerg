package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg"
)

type followersParser struct {
	name string
}

func NewFollowersParser(name string) *followersParser {
	return &followersParser{
		name: name,
	}
}

func (p *followersParser) Name() string {
	return p.name
}

func (p *followersParser) Parse(rsp *zerg.Response) ([]*zerg.Item, []*zerg.Request) {
	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		log.Error(err)
		return nil, nil
	}
	var reqs []*zerg.Request
	doc.Find("#Profile-following .ContentItem .Popover a").Each(
		func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			href = "http:" + href
			req, _ := zerg.NewRequest(href+"/followers?page=1", "followers")
			reqs = append(reqs, req)
			req, _ = zerg.NewRequest(href+"/following?page=1", "followers")
			reqs = append(reqs, req)
			req, _ = zerg.NewRequest(href+"/following/questions?page=1", "questions")
			reqs = append(reqs, req)
		})
	if len(reqs) > 0 {
		tokens := strings.Split(rsp.RawURL, "?page=")
		if len(tokens) == 2 {
			if page, err := strconv.Atoi(tokens[1]); err == nil {
				url := fmt.Sprintf("%s?page=%d", tokens[0], page+1)
				req, _ := zerg.NewRequest(url, "followers")
				reqs = append(reqs, req)
			}
		}
	}
	return nil, reqs
}
