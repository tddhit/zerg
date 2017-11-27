package jobbole

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/spider"
)

type Parser struct {
	reg *regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		reg: regexp.MustCompile(`(\d+).*`),
	}
}

func (p *Parser) Parse(rsp *engine.Response) (item *engine.Item, reqs []*engine.Request) {
	item.Dict["url"] = rsp.RawUrl
	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	title := doc.Find(".entry-header h1").Text()
	if title != "" {
		p.parseDetailPage(doc, item, reqs)
	} else {
		doc.Find("#archive .post-thumb a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := http.NewRequest("GET", href, nil)
			ireq := &engine.Request{
				req,
				RawUrl: href,
			}
			reqs = append(reqs, req)
		})
		doc.Find(".next.page-numbers").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req, _ := http.NewRequest("GET", href, nil)
			ireq := &engine.Request{
				req,
				RawUrl: href,
			}
			reqs = append(reqs, req)
		})
	}
}

func (this *Processer) ProcessDetail(doc *goquery.Document, item *engine.Item, reqs []*engine.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	title := doc.Find(".entry-header h1").Text()
	createDate := doc.Find(".entry-meta-hide-on-mobile").Text()
	createDate = strings.Replace(createDate, "·", " ", -1)
	createDate = strings.Split(strings.TrimSpace(createDate), " ")[0]
	voteNum := doc.Find(".vote-post-up h10").Text()
	bookmarkNum := doc.Find(".bookmark-btn").Text()
	tmpBookmarkNum := this.reg.FindStringSubmatch(bookmarkNum)
	if len(tmpBookmarkNum) > 1 {
		bookmarkNum = tmpBookmarkNum[1]
	} else {
		bookmarkNum = ""
	}
	commentNum := doc.Find("a[href='#article-comment']").Text()
	tmpCommentNum := this.reg.FindStringSubmatch(commentNum)
	if len(tmpCommentNum) > 1 {
		commentNum = tmpCommentNum[1]
	} else {
		commentNum = ""
	}
	tags := make([]string, 0)
	doc.Find(".entry-meta-hide-on-mobile a[href^='http://blog.jobbole.com/']").Each(func(i int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})
	content := doc.Find(".entry p").Text()
	content = strings.Join(strings.Fields(content), " ")
	item.Dict["url"] = url
	item.Dict["title"] = title
	item.Dict["createDate"] = createDate
	item.Dict["voteNum"] = voteNum
	item.Dict["bookmarkNum"] = bookmarkNum
	item.Dict["commentNum"] = commentNum
	item.Dict["tags"] = strings.Join(tags, ",")
	item.Dict["content"] = content
}
