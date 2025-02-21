package stdgz

import (
	"compress/gzip"
	"context"
	"encoding/binary"
	"fmt"
	"io"

	gs "github.com/takanoriyanagitani/go-gzip-stat"
	s2 "github.com/takanoriyanagitani/go-gzip-stat/gzip/stat"
	. "github.com/takanoriyanagitani/go-gzip-stat/util"
)

type GzReader struct {
	*gzip.Reader
}

func (g GzReader) PopulateStatus(s *gs.GzipStatus) {
	s.Comment = g.Reader.Header.Comment
	s.Modified = g.Reader.Header.ModTime
	s.Name = g.Reader.Header.Name
	s.OS = g.Reader.Header.OS
}

func GzipFileLikeToStatNew() s2.GzipFileToStat {
	var buf gs.GzipStatus
	var msize [4]byte
	var grdr *gzip.Reader

	return func(gzfile s2.NamedFile) IO[gs.GzipStatus] {
		return func(_ context.Context) (gs.GzipStatus, error) {
			size, e := gzfile.Seek(0, io.SeekEnd)
			if nil != e {
				return gs.GzipStatus{}, fmt.Errorf("%w: unable to seek", e)
			}

			buf.Size = uint64(size)
			_, e = gzfile.Seek(-4, io.SeekEnd)
			if nil != e {
				return gs.GzipStatus{}, fmt.Errorf("%w: unable to get size", e)
			}

			_, e = io.ReadFull(gzfile, msize[:])
			if nil != e {
				return gs.GzipStatus{}, fmt.Errorf("%w: unable to read size", e)
			}

			var usize uint32 = binary.LittleEndian.Uint32(msize[:])

			buf.OriginalSizeMod32 = usize

			_, e = gzfile.Seek(0, io.SeekStart)
			if nil != e {
				return gs.GzipStatus{}, fmt.Errorf("%w: unable to seek", e)
			}

			switch grdr {
			case nil:
				grdr, e = gzip.NewReader(gzfile)
			default:
				e = grdr.Reset(gzfile)
			}
			if nil != e {
				return gs.GzipStatus{}, fmt.Errorf("%w: invalid gzip", e)
			}

			GzReader{grdr}.PopulateStatus(&buf)

			if 0 == len(buf.Name) {
				buf.Name = gzfile.Filename
			}

			return buf, nil
		}
	}
}
