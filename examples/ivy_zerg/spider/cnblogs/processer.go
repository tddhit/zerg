package jobbole

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/go_spider/core/common/page"
	"github.com/tddhit/go_spider/core/common/request"
)

const timeFormat = "2006-01-02 15:04"

type Processer struct {
}

func NewProcesser() *Processer {
	return &Processer{}
}

func (this *Processer) ProcessDetail(p *page.Page) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	query := p.GetHtmlParser()
	url := p.GetRequest().GetUrl()
	title := query.Find("#cb_post_title_url").Text()
	createDate := query.Find("#post-date").Text()
	ts, _ := time.Parse(timeFormat, createDate)
	createDate = strconv.FormatUint(uint64(ts.Unix()), 10)
	//voteNum := query.Find(".diggnum").Text()
	//commentNum := query.Find("#post_comment_count").Text()
	//viewNum := query.Find("#post_view_count").Text()
	//tags := make([]string, 0)
	//query.Find("#EntryTag a").Each(func(i int, s *goquery.Selection) {
	//	tags = append(tags, s.Text())
	//})
	content := query.Find("#cnblogs_post_body p,h2,h3").Text()
	content = strings.Join(strings.Fields(content), " ")
	p.AddField("url", url)
	p.AddField("title", title)
	p.AddField("createDate", createDate)
	p.AddField("voteNum", voteNum)
	p.AddField("bookmarkNum", "0")
	p.AddField("viewNum", viewNum)
	p.AddField("commentNum", commentNum)
	p.AddField("tags", strings.Join(tags, ","))
	p.AddField("content", content)
}

func (this *Processer) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}
	query := p.GetHtmlParser()
	title := query.Find("#cb_post_title_url").Text()
	if title != "" {
		this.ProcessDetail(p)
		return
	} else {
		query.Find("#post_list .post_item_body h3 a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req := request.NewRequest(href, "html", "jobbole", "GET", "", nil, nil, nil, nil)
			p.AddTargetRequestWithParams(req)
		})
		query.Find("#pager_bottom .pager :last-child").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			req := request.NewRequest(href, "html", "jobbole", "GET", "", nil, nil, nil, nil)
			p.AddTargetRequestWithParams(req)
		})
	}
}

func (this *Processer) Finish() {
	log.Printf("done！ \r\n")
}
