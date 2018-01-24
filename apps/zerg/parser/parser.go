package parser

import (
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/types"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Parser struct {
	name        string
	typ         string
	cssSelector string
	writer      string
	parser      string
}

func NewParser(name, typ, cssSelector, writer, parser string) *Parser {
	return &Parser{
		name:        name,
		typ:         typ,
		cssSelector: cssSelector,
		writer:      writer,
		parser:      parser,
	}
}

func (p *Parser) Name() string {
	return p.name
}

func (p *Parser) Parse(rsp *types.Response) ([]*types.Item, []*types.Request) {
	var items []*types.Item
	var reqs []*types.Request
	if p.typ == "href" {
		reqs = p.parseHref(rsp)
	} else if p.typ == "text" {
		items = p.parseText(rsp)
	}
	return items, reqs
}

func (p *Parser) parseHref(rsp *types.Response) []*types.Request {
	var reqs []*types.Request
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	if doc == nil {
		return nil
	}
	doc.Find(p.cssSelector).Each(func(i int, contentSelection *goquery.Selection) {
		href, _ := contentSelection.Attr("href")
		log.Debug(href)
		req, _ := types.NewRequest(href, p.parser)
		reqs = append(reqs, req)
	})
	return reqs
}

func (p *Parser) parseText(rsp *types.Response) []*types.Item {
	var items []*types.Item
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	if doc == nil {
		return nil
	}
	doc.Find(p.cssSelector).Each(func(i int, contentSelection *goquery.Selection) {
		var content string
		paragraph := contentSelection.Text()
		detector := chardet.NewTextDetector()
		result, err := detector.DetectBest([]byte(paragraph))
		if err == nil {
			if result.Charset != "UTF-8" && result.Confidence > 80 {
				sr := strings.NewReader(paragraph)
				tr := transform.NewReader(sr, simplifiedchinese.GB18030.NewDecoder())
				b, err := ioutil.ReadAll(tr)
				if err == nil {
					content += strings.Join(strings.Fields(string(b)), " ")
				}
			} else if result.Charset == "UTF-8" && result.Confidence > 80 {
				content += strings.Join(strings.Fields(paragraph), " ")
			}
		}
		if content != "" {
			item := types.NewItem(p.writer)
			item.Dict["url"] = rsp.RawURL
			item.Dict["content"] = content
			items = append(items, item)
		}
	})
	return items
}
