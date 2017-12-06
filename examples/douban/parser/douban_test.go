package parser

import (
	"log"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

const timeFormat = "2006-01-02 15:04"

func TestExample(t *testing.T) {
	{
		query, err := goquery.NewDocument("https://movie.douban.com/review/8965655/")
		if err != nil {
			log.Fatal(err)
		}
		title := query.Find("h1 span").Text()
		log.Println(title)
		content := query.Find(".main-bd p").Text()
		log.Println(content)
	}
	{
		query, err := goquery.NewDocument("https://movie.douban.com/review/best/")
		if err != nil {
			log.Fatal(err)
		}
		query.Find(".main-bd h2 a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			log.Println(href)
		})
	}
	{
		query, err := goquery.NewDocument("https://movie.douban.com/review/best/")
		if err != nil {
			log.Fatal(err)
		}
		query.Find(".next a").Each(func(i int, contentSelection *goquery.Selection) {
			href, _ := contentSelection.Attr("href")
			log.Println(href)
		})
	}
}
