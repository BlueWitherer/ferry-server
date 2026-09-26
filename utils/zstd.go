package utils

import (
	"sync"

	"github.com/klauspost/compress/zstd"
	"github.com/samber/mo"
)

var (
	zstdEncoder *zstd.Encoder
	zstdDecoder *zstd.Decoder
	zstdOnce    sync.Once
)

func initZstd() {
	zstdEncoder, _ = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBetterCompression))
	zstdDecoder, _ = zstd.NewReader(nil)
}

func CompressZstd(data []byte) []byte {
	zstdOnce.Do(initZstd)
	return zstdEncoder.EncodeAll(data, make([]byte, 0, len(data)))
}

func DecompressZstd(data []byte) mo.Result[[]byte] {
	zstdOnce.Do(initZstd)

	if zstdDecoder == nil {
		return mo.Errf[[]byte]("zstd decoder failed to initialize")
	}

	bytes, err := zstdDecoder.DecodeAll(data, nil)
	if err != nil {
		return mo.Err[[]byte](err)
	}

	return mo.Ok(bytes)
}
