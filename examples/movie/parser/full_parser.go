package parser

import (
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"github.com/tddhit/zerg/types"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type FullParser struct {
	name string
}

func NewFullParser(name string) *FullParser {
	return &FullParser{
		name: name,
	}
}

func (p *FullParser) Name() string {
	return p.name
}

func (p *FullParser) Parse(rsp *types.Response) ([]*types.Item, []*types.Request) {
	var items []*types.Item
	//item.Dict["url"] = rsp.RawURL
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	doc.Find("p").Each(func(i int, contentSelection *goquery.Selection) {
		item := types.NewItem("full")
		var content string
		paragraph := contentSelection.Text()
		detector := chardet.NewTextDetector()
		result, err := detector.DetectBest([]byte(paragraph))
		if err == nil {
			if result.Charset != "UTF-8" && result.Language == "zh" && result.Confidence > 80 {
				sr := strings.NewReader(paragraph)
				tr := transform.NewReader(sr, simplifiedchinese.GB18030.NewDecoder())
				b, err := ioutil.ReadAll(tr)
				if err == nil {
					content += strings.Join(strings.Fields(string(b)), " ")
				}
			} else if result.Charset == "UTF-8" && result.Language == "zh" && result.Confidence > 80 {
				content += strings.Join(strings.Fields(paragraph), " ")
			}
		}
		if content != "" {
			item.Dict["content"] = content
			items = append(items, item)
		}
	})
	return items, nil
}
