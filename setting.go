package uRouter

var (
	defaultJsonCodec Codec = new(stdJsonCodec)

	TextHeader = NewHeaderCodec(TextLengthEncoding, defaultJsonCodec, func() Header {
		return &MapHeader{}
	})

	BinaryHeader = NewHeaderCodec(BinaryLengthEncoding, defaultJsonCodec, func() Header {
		return &MapHeader{}
	})
)

func SetJsonCodec(codec Codec) {
	defaultJsonCodec = codec

	TextHeader = NewHeaderCodec(TextLengthEncoding, defaultJsonCodec, func() Header {
		return &MapHeader{}
	})

	BinaryHeader = NewHeaderCodec(BinaryLengthEncoding, defaultJsonCodec, func() Header {
		return &MapHeader{}
	})
}
