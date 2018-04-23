package parser

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/tools/log"
)

func TestParser(t *testing.T) {
	request, err := http.NewRequest("GET", "http://music.163.com/artist/album?id=1011189", nil)
	if err != nil {
		log.Fatal(err)
	}
	//request.Header.Add("Referer", "http://music.163.com/")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Debug(string(body))
	//os.Exit(1)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	b, _ := doc.Html()
	log.Debug(b)
	doc.Find("#m-song-module li p a").Each(func(i int, contentSelection *goquery.Selection) {
		href, _ := contentSelection.Attr("href")
		log.Debug(href)
		//log.Info(ALBUM_ROOTURL+v[1], "NeteaseMusicParser")
	})
	doc.Find(".u-page .znxt").Each(func(i int, contentSelection *goquery.Selection) {
		href, _ := contentSelection.Attr("href")
		log.Debug(href)
	})
}
