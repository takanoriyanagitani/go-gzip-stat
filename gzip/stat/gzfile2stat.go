package stat

import (
	"bufio"
	"context"
	"io"
	"iter"
	"log"
	"os"

	gs "github.com/takanoriyanagitani/go-gzip-stat"
	. "github.com/takanoriyanagitani/go-gzip-stat/util"
)

type NamedFile struct {
	Filename string
	io.ReadSeekCloser
}

type GzipFileToStat func(NamedFile) IO[gs.GzipStatus]

type NameToFileLike func(string) IO[io.ReadSeekCloser]

type StatWriter func(gs.GzipStatus) IO[Void]

type NamesToStats func(iter.Seq[string]) IO[iter.Seq2[gs.GzipStatus, error]]

func (g GzipFileToStat) ToNamesToStats(n2f NameToFileLike) NamesToStats {
	return func(names iter.Seq[string]) IO[iter.Seq2[gs.GzipStatus, error]] {
		return func(
			ctx context.Context,
		) (iter.Seq2[gs.GzipStatus, error], error) {
			return func(yield func(gs.GzipStatus, error) bool) {
				var empty gs.GzipStatus
				for name := range names {
					select {
					case <-ctx.Done():
						yield(empty, ctx.Err())
						return
					default:
					}

					rsc, e := n2f(name)(ctx)
					if nil != e {
						yield(empty, e)
						return
					}

					stat, e := func() (gs.GzipStatus, error) {
						defer rsc.Close()

						named := NamedFile{
							Filename:       name,
							ReadSeekCloser: rsc,
						}

						return g(named)(ctx)
					}()
					if !yield(stat, e) {
						return
					}
				}
			}, nil
		}
	}
}

type NamesToStatsToWriter func(iter.Seq[string]) IO[Void]

func (n NamesToStats) ToNamesToStatsToWriter(
	wtr StatWriter,
) NamesToStatsToWriter {
	return func(names iter.Seq[string]) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			stats, e := n(names)(ctx)
			if nil != e {
				return Empty, e
			}

			for stat, e := range stats {
				if nil != e {
					return Empty, e
				}

				_, e := wtr(stat)(ctx)
				if nil != e {
					return Empty, e
				}
			}

			return Empty, nil
		}
	}
}

func (n NamesToStatsToWriter) StdinToNamesToWriter() IO[Void] {
	var s *bufio.Scanner = bufio.NewScanner(os.Stdin)
	var i iter.Seq[string] = func(
		yield func(string) bool,
	) {
		for s.Scan() {
			var filename string = s.Text()
			if !yield(filename) {
				return
			}
		}

		e := s.Err()
		if nil != e {
			log.Printf("Error while scanning filenames: %v\n", e)
		}
	}

	return n(i)
}
