package main

import (
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/zerg/engine"
	"github.com/tddhit/zerg/examples/music/internal/header"
	"github.com/tddhit/zerg/examples/music/parser"
	"github.com/tddhit/zerg/examples/music/queuer"
	"github.com/tddhit/zerg/examples/music/writer"
)

func main() {
	neteaseArtistParser := parser.NewNeteaseArtistParser("NeteaseArtistParser")
	neteaseAlbumParser := parser.NewNeteaseAlbumParser("NeteaseAlbumParser")
	neteaseMusicParser := parser.NewNeteaseMusicParser("NeteaseMusicParser")
	neteaseMusicWriter := writer.NewFileWriter("NeteaseMusicWriter", "/home/tdd/go/src/github.com/tddhit/zerg/examples/music/data/netease.txt")

	engine := engine.NewEngine(engine.Option{LogLevel: log.INFO})
	engine.AddParser(neteaseArtistParser).AddParser(neteaseAlbumParser)
	engine.AddParser(neteaseMusicParser).AddWriter(neteaseMusicWriter)
	for i := 65; i <= 90; i++ {
		engine.AddSeed("http://music.163.com/discover/artist/cat?id=1001&initial="+strconv.Itoa(i), "NeteaseArtistParser", "", header.Header)
		engine.AddSeed("http://music.163.com/discover/artist/cat?id=1002&initial="+strconv.Itoa(i), "NeteaseArtistParser", "", header.Header)
		engine.AddSeed("http://music.163.com/discover/artist/cat?id=1003&initial="+strconv.Itoa(i), "NeteaseArtistParser", "", header.Header)
	}
	log.Info(header.Header)
	go func() {
		log.Debug(http.ListenAndServe("localhost:6060", nil))
	}()
	//engine.AddSeed("http://music.163.com/artist/album?id=8985", "NeteaseAlbumParser", "", header.Header)
	engine.SetSchedulerPolicy(queuer.NewDefaultQueuer())
	engine.Go()
}
