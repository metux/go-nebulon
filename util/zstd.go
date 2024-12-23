package util

import (
	"time"

	"github.com/klauspost/compress/zstd"
)

// configuration section
var (
	//	ZipLevel = zstd.SpeedFastest
	ZipLevel = zstd.SpeedBestCompression
	TraceZip = false
)

var (
	zipWriter, zipErrW = zstd.NewWriter(nil, zstd.WithEncoderLevel(ZipLevel))
	zipReader, zipErrR = zstd.NewReader(nil)
)

func ZipEncode(data []byte) []byte {
	if TraceZip {
		defer TimeTrack(time.Now(), "ZipEncode")
	}
	return zipWriter.EncodeAll(data, make([]byte, 0, len(data)))
}

func ZipDecode(data []byte) ([]byte, error) {
	if TraceZip {
		defer TimeTrack(time.Now(), "ZipDecode")
	}
	return zipReader.DecodeAll(data, make([]byte, 0, len(data)))
}
