package constant

/*
定义了内部使用的一些常量
*/

const (
	XPath       = "X-Path"
	XRealIP     = "X-Real-Ip"
	ContentType = "Content-Type"
)

const (
	MimeJson   = "application/json; charset=utf-8"
	MimeStream = "application/octet-stream"
)

const (
	KiB           = 1024
	BufferLeveL1  = KiB
	BufferLeveL4  = 4 * KiB
	BufferLeveL8  = 8 * KiB
	BufferLeveL16 = 16 * KiB
)

const (
	HttpHeaderNumber = iota
	MapHeaderNumber
)
