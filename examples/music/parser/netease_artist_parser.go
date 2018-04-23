package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/examples/music/internal/header"
	"github.com/tddhit/zerg/types"
)

const ARTIST_ROOTURL = "http://music.163.com/artist/album?"

type NeteaseArtistParser struct {
	name string
}

func NewNeteaseArtistParser(name string) *NeteaseArtistParser {
	return &NeteaseArtistParser{
		name: name,
	}
}

func (p *NeteaseArtistParser) Name() string {
	return p.name
}

func (p *NeteaseArtistParser) Parse(rsp *types.Response) ([]*types.Item, []*types.Request) {
	var reqs []*types.Request
	doc, _ := goquery.NewDocumentFromReader(rsp.Body)
	if doc == nil {
		return nil, nil
	}
	doc.Find(".m-sgerlist li .nm").Each(func(i int, contentSelection *goquery.Selection) {
		href, _ := contentSelection.Attr("href")
		v := strings.Split(href, "?")
		if len(v) != 2 {
			log.Error(href)
			return
		}
		req, _ := types.NewRequest(ARTIST_ROOTURL+v[1], "NeteaseAlbumParser", "", header.Header)
		reqs = append(reqs, req)
		log.Info(ARTIST_ROOTURL+v[1], "NeteaseAlbumParser")
	})
	return nil, reqs
}
