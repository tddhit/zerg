package goquery_test

import (
	"log"
	//"strconv"
	//"strings"
	"testing"
	//"time"

	"github.com/PuerkitoBio/goquery"
)

const timeFormat = "2006-01-02 15:04"

func TestExample(t *testing.T) {
	query, err := goquery.NewDocument("http://www.cnblogs.com/weiqinl/p/7886049.html")
	if err != nil {
		log.Fatal(err)
	}

	query.Find("#post_list .post_item_body h3 a").Each(func(i int, contentSelection *goquery.Selection) {
		href, _ := contentSelection.Attr("href")
		log.Println(href)
	})

	query.Find("#pager_bottom .pager :last-child").Each(func(i int, contentSelection *goquery.Selection) {
		href, _ := contentSelection.Attr("href")
		log.Println(href)
	})

	query.Find("#cb_post_title_url").Each(func(i int, contentSelection *goquery.Selection) {
		title := contentSelection.Text()
		log.Println(title)
	})
	//title := query.Find("#cb_post_title_url").Text()
	createDate := query.Find(".postDesc").Text()
	//ts, _ := time.Parse(timeFormat, createDate)
	//createDate = strconv.FormatUint(uint64(ts.Unix()), 10)
	//voteNum := query.Find(".diggnum").Text()
	//commentNum := query.Find("#post_comment_count").Text()
	viewNum := query.Find("#post_view_count").Text()
	//tags := make([]string, 0)
	//query.Find("#EntryTag a").Each(func(i int, s *goquery.Selection) {
	//	tags = append(tags, s.Text())
	//})
	//content := query.Find("#cnblogs_post_body p,h2,h3").Text()
	//content = strings.Join(strings.Fields(content), " ")
	//log.Println("title", title)
	log.Println("createDate", createDate)
	//log.Println("voteNum", voteNum)
	//log.Println("bookmarkNum", "0")
	log.Println("viewNum", viewNum)
	//log.Println("commentNum", commentNum)
	//log.Println("tags", strings.Join(tags, ","))
	//log.Println("content", content)
}
