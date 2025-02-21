package jwriter

import (
	"context"
	"encoding/json"
	"io"

	gs "github.com/takanoriyanagitani/go-gzip-stat"
	s2 "github.com/takanoriyanagitani/go-gzip-stat/gzip/stat"
	. "github.com/takanoriyanagitani/go-gzip-stat/util"
)

func StatWriterJsonNew(wtr io.Writer) s2.StatWriter {
	var enc *json.Encoder = json.NewEncoder(wtr)
	return func(s gs.GzipStatus) IO[Void] {
		return func(_ context.Context) (Void, error) {
			return Empty, enc.Encode(&s)
		}
	}
}
