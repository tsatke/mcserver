package chunk

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
)

//go:generate stringer -linecomment -type=CompressionType

type CompressionType byte

const (
	CompressionGZip         CompressionType = iota + 1 // GZip
	CompressionZlib                                    // Zlib
	CompressionUncompressed                            // uncompressed
)

var (
	Decompressor = [4]func(io.Reader) (io.Reader, error){
		func(_ io.Reader) (io.Reader, error) {
			return nil, fmt.Errorf("unknown compression 0x00")
		},
		func(in io.Reader) (io.Reader, error) {
			return gzip.NewReader(in)
		},
		func(in io.Reader) (io.Reader, error) {
			return zlib.NewReader(in)
		},
		func(in io.Reader) (io.Reader, error) {
			return in, nil
		},
	}
)
