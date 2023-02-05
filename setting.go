package uRouter

var (
	defaultJsonCodec Codec = new(stdJsonCodec)
)

func SetJsonCodec(codec Codec) {
	defaultJsonCodec = codec
}
