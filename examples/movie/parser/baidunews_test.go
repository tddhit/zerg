package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const timeFormat = "2006-01-02 15:04"

func TestExample(t *testing.T) {
	{
		doc, err := goquery.NewDocument("http://news.163.com/16/0511/07/BMP2VFA600014U9R.html")
		if err != nil {
			log.Fatal(err)
		}
		doc.Find("p").Each(func(i int, contentSelection *goquery.Selection) {
			content := contentSelection.Text()
			detector := chardet.NewTextDetector()
			result, err := detector.DetectBest([]byte(content))
			if err == nil {
				fmt.Printf(
					"Detected charset is %s, language is %s\n",
					result.Charset,
					result.Language)
			}
			if result.Charset != "UTF-8" {
				sr := strings.NewReader(content)
				tr := transform.NewReader(sr, simplifiedchinese.GB18030.NewDecoder())
				b, err := ioutil.ReadAll(tr)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(b))
			} else {
				fmt.Println(content)
			}
		})
	}
}
