package main

import (
	"context"
	"io"
	"log"
	"os"

	s2 "github.com/takanoriyanagitani/go-gzip-stat/gzip/stat"
	js "github.com/takanoriyanagitani/go-gzip-stat/gzip/stat/stat2writer/json/std"
	s3 "github.com/takanoriyanagitani/go-gzip-stat/gzip/stat/std"
	. "github.com/takanoriyanagitani/go-gzip-stat/util"
)

var stat2json2stdout s2.StatWriter = js.StatWriterJsonNew(os.Stdout)

var name2file s2.NameToFileLike = Lift(
	func(filename string) (io.ReadSeekCloser, error) { return os.Open(filename) },
)

var gzfile2stat s2.GzipFileToStat = s3.GzipFileLikeToStatNew()

var names2stats s2.NamesToStats = gzfile2stat.ToNamesToStats(name2file)

var names2stats2json2stdout s2.NamesToStatsToWriter = names2stats.
	ToNamesToStatsToWriter(stat2json2stdout)

var stdin2gzfilenames2stats2json2writer IO[Void] = names2stats2json2stdout.
	StdinToNamesToWriter()

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return stdin2gzfilenames2stats2json2writer(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
