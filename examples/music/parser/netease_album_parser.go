package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/examples/music/internal/header"
	"github.com/tddhit/zerg/types"
)

const ALBUM_ROOTURL = "http://music.163.com/album?"

type NeteaseAlbumParser struct {
	name string
}

func NewNeteaseAlbumParser(name string) *NeteaseAlbumParser {
	return &NeteaseAlbumParser{
		name: name,
	}
}

func (p *NeteaseAlbumParser) Name() string {
	return p.name
}

func (p *NeteaseAlbumParser) Parse(rsp *types.Response) ([]*types.Item, []*types.Request) {
	log.Info(rsp.RawURL)
	var reqs []*types.Request
	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if doc == nil {
		log.Error(rsp.RawURL, err)
		return nil, nil
	}
	doc.Find("#m-song-module li p a").Each(func(i int, contentSelection *goquery.Selection) {
		href, _ := contentSelection.Attr("href")
		v := strings.Split(href, "?")
		if len(v) != 2 {
			log.Error(href)
			return
		}
		req, _ := types.NewRequest(ALBUM_ROOTURL+v[1], "NeteaseMusicParser", "", header.Header)
		reqs = append(reqs, req)
		//log.Info(ALBUM_ROOTURL+v[1], "NeteaseMusicParser")
	})
	doc.Find(".u-page .znxt").Each(func(i int, contentSelection *goquery.Selection) {
		href, _ := contentSelection.Attr("href")
		v := strings.Split(href, "?")
		if len(v) != 2 {
			log.Error(href)
			return
		}
		req, _ := types.NewRequest(ARTIST_ROOTURL+v[1], "NeteaseAlbumParser", "", header.Header)
		reqs = append(reqs, req)
		//log.Info(ARTIST_ROOTURL+v[1], "NeteaseAlbumParser")
	})
	return nil, reqs
}
