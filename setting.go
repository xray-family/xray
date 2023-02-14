package uRouter

import (
	"github.com/lxzan/uRouter/codec"
)

var (
	defaultJsonCodec codec.Codec = codec.StdJsonCodec

	defaultBufferPool BufferPoolInterface = newBufferPool()

	defaultLogger LoggerInterface = new(logger)

	defaultHeaderPool = newHeaderPool()

	MapHeaderTemplate = new(MapHeader) // readonly
)

// SetJsonCodec 设置JSON编码器
// 如果Header不使用JSON编码, 调用SetJsonCodec后需要开发者重新配置HeaderCodec
// if the Header does not use JSON encoding, the developer needs to reconfigure the HeaderCodec after calling SetJsonCodec
func SetJsonCodec(codec codec.Codec) {
	defaultJsonCodec = codec

	TextMapHeader = NewHeaderCodec(MapHeaderTemplate, codec).setLengthBytes(textLengthEncoding)

	BinaryMapHeader = NewHeaderCodec(MapHeaderTemplate, codec).setLengthBytes(binaryLengthEncoding)
}

// SetBufferPool 设置buffer池
func SetBufferPool(p BufferPoolInterface) {
	defaultBufferPool = p
}

// SetLogger 设置日志工具
func SetLogger(l LoggerInterface) {
	defaultLogger = l
}

// BufferPool 获取默认缓冲池
func BufferPool() BufferPoolInterface {
	return defaultBufferPool
}

// Logger 获取默认日志工具
func Logger() LoggerInterface {
	return defaultLogger
}

// JsonCodec 获取默认JSON编码器
func JsonCodec() codec.Codec {
	return defaultJsonCodec
}
