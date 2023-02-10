package internal

type UError string

func (c UError) Error() string {
	return string(c)
}

const (
	ErrHeaderTooLarge    = UError("header size too large")
	ErrParseHeader       = UError("parse header error")
	ErrIOBytesUnexpected = UError("io bytes unexpected")
	ErrNilRequestBody    = UError("nil request body")
)
