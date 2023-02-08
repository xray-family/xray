package uRouter

var (
	defaultJsonCodec Codec = new(stdJsonCodec)

	defaultGenerator = func() Header { return &MapHeader{} }

	defaultBufferPool BufferPool = newBufferPool()

	defaultLogger Logger = new(logger)
)

var (
	StdJsonCodec Codec = new(stdJsonCodec)

	// TextHeader 文本类型头部编码, 4字节, 最大长度=9999
	// text type header code, 4 bytes, max length = 9999
	TextHeader *HeaderCodec

	// BinaryHeader 二进制类型头部编码, 2字节, 最大长度=65535
	// binary type header code, 2 bytes, max length = 65535
	BinaryHeader *HeaderCodec
)

func init() {
	TextHeader = NewHeaderCodec(defaultJsonCodec, defaultGenerator).setLengthBytes(textLengthEncoding)

	BinaryHeader = NewHeaderCodec(defaultJsonCodec, defaultGenerator).setLengthBytes(binaryLengthEncoding)
}

// SetJsonCodec 设置JSON编码器
// 如果Header不使用JSON编码, 调用SetJsonCodec后需要开发者重新配置HeaderCodec
// if the Header does not use JSON encoding, the developer needs to reconfigure the HeaderCodec after calling SetJsonCodec
func SetJsonCodec(codec Codec) {
	defaultJsonCodec = codec

	TextHeader = NewHeaderCodec(defaultJsonCodec, defaultGenerator).setLengthBytes(textLengthEncoding)

	BinaryHeader = NewHeaderCodec(defaultJsonCodec, defaultGenerator).setLengthBytes(binaryLengthEncoding)
}

// SetBufferPool 设置buffer池
func SetBufferPool(p BufferPool) {
	defaultBufferPool = p
}

// SetLogger 设置日志工具
func SetLogger(l Logger) {
	defaultLogger = l
}
