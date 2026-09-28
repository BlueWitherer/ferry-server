package utils

import (
	"encoding/json"
	"sync"

	"github.com/klauspost/compress/zstd"
	"github.com/samber/mo"
)

var (
	zstdEncoder *zstd.Encoder
	zstdDecoder *zstd.Decoder

	zstdOnce sync.Once
)

func initZstd() {
	zstdEncoder, _ = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBetterCompression))
	zstdDecoder, _ = zstd.NewReader(nil)
}

func ZstdCompress(data []byte) []byte {
	zstdOnce.Do(initZstd)
	return zstdEncoder.EncodeAll(data, make([]byte, 0, len(data)))
}

func ZstdDecompress(data []byte) mo.Result[[]byte] {
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

func ZstdCompressJSON[T any](v T) mo.Result[[]byte] {
	b, err := json.Marshal(v)
	if err != nil {
		return mo.Err[[]byte](err)
	}

	return mo.Ok(ZstdCompress(b))
}

func ZstdDecompressJSON[T any](data []byte) mo.Result[T] {
	decompressedRes := ZstdDecompress(data)
	if decompressedRes.IsError() {
		return mo.Err[T](decompressedRes.Error())
	}

	var out T
	if err := json.Unmarshal(decompressedRes.MustGet(), &out); err != nil {
		return mo.Err[T](err)
	}

	return mo.Ok(out)
}
